package sciencedirect

import (
	"context"
	"github.com/antchfx/htmlquery"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	"web_spider/pkg"
)

type ScienceDirect struct {
	log *log.Helper

	chrome  *pkg.ChromePool
	workDir string
	page    *rod.Page

	ctx        context.Context
	cancel     context.CancelFunc
	maxWorkers int
	wg         sync.WaitGroup
	saveChan   chan saveTask
	domain     string
}

type saveTask struct {
	targetUrl string
	reqId     proto.NetworkRequestID
}

func NewScienceDirect(chrome *pkg.ChromePool, logger log.Logger) (*ScienceDirect, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	workDir := path.Join(homeDir, "sciencedirect")
	err = os.MkdirAll(workDir, 0755)
	if err != nil {
		ctxCancel()
		return nil, err
	}

	sd := &ScienceDirect{
		cancel:     ctxCancel,
		ctx:        ctx,
		log:        log.NewHelper(logger),
		chrome:     chrome,
		maxWorkers: 5,
		workDir:    workDir,
		saveChan:   make(chan saveTask),
		domain:     "https://www.sciencedirect.com",
	}

	for i := 0; i < sd.maxWorkers; i++ {
		sd.wg.Add(1)
		go sd.saveWorker()
	}

	return sd, nil
}

func (s *ScienceDirect) Close() error {
	s.log.Infof("Closing ScienceDirect")
	s.cancel()

	s.log.Infof("start wait for all workers to finish")
	s.wg.Wait()
	s.log.Infof("all workers finished")

	return nil
}

func (s *ScienceDirect) saveWorker() {
	defer s.wg.Done()

	for {
		select {
		case task := <-s.saveChan:
			if err := s.saveDir(task.targetUrl, task.reqId); err != nil {
				s.log.Errorf("Failed to save response: %v", err)
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *ScienceDirect) List() error {
	page, err := s.chrome.Browser.Page(proto.TargetCreateTarget{
		URL: s.domain,
	})

	if err != nil {
		return err
	}
	s.page = page

	go page.EachEvent(func(e *proto.NetworkResponseReceived) {
		respUrl := e.Response.URL

		if strings.HasPrefix(respUrl, s.domain) {
			s.log.Debugf("new resp url: %v", respUrl)
		}

		if strings.HasPrefix(respUrl, "https://www.sciencedirect.com/search/api") {
			s.log.Infof("sciencedirect api url: %s", respUrl)

			select {
			case s.saveChan <- saveTask{targetUrl: respUrl, reqId: e.RequestID}:
			case <-s.ctx.Done():
				s.log.Infof("context done, stop saving")
				return
			}
		}
	})()

	<-s.ctx.Done()
	return nil
}

func (s *ScienceDirect) saveDir(targetUrl string, reqId proto.NetworkRequestID) (err error) {
	s.log.Infof("start new task for target url: %s", targetUrl)

	time.Sleep(5 * time.Second) // wait for response

	m := proto.NetworkGetResponseBody{RequestID: reqId}
	resp, err := m.Call(s.page)
	if err != nil {
		s.log.Errorf("page.call err: %+v", err)
		return err
	}

	up, _ := url.Parse(targetUrl)
	qs := up.Query().Get("qs")

	qsDir := path.Join(s.workDir, qs)
	_ = os.MkdirAll(qsDir, 0755)

	result := gjson.Parse(resp.Body)

	searchResults := result.Get("searchResults").Array()

	s.log.Infof("search results size: %d", len(searchResults))

	for _, sr := range searchResults {
		pii := sr.Get("pii").String()
		resultPath := path.Join(qsDir, pii+".json")

		if pkg.FileExists(resultPath) {
			continue
		}

		err = os.WriteFile(resultPath, []byte(sr.String()), 0755)
		if err != nil {
			s.log.Errorf("save to local dir err: %+v", err)
			return err
		}
	}
	return
}

func (s *ScienceDirect) Detail(workDir string) error {
	if workDir == "" {
		s.log.Infof("workDir is empty")
		return nil
	}
	page, err := s.chrome.Browser.Page(proto.TargetCreateTarget{
		URL: s.domain,
	})

	if err != nil {
		return err
	}
	s.page = page

	urls, err := s.parseAllLink(workDir)
	s.log.Infof("for dir: %s, size: %d", workDir, len(urls))

	if err != nil {
		s.log.Errorf("process err: %+v", err)
		return err
	}

	for _, urlItem := range urls {
		err = s.process(urlItem)
		if err != nil {
			s.log.Errorf("process err: %+v", err)
			return err
		}
	}
	return nil
}

func (s *ScienceDirect) parseAllLink(workDir string) ([]string, error) {
	entries, err := os.ReadDir(workDir)
	if err != nil {
		s.log.Errorf("ReadDir err: %+v", err)
		return nil, err
	}

	var results []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".json") {
			abstractFile := strings.ReplaceAll(entry.Name(), ".json", ".abstract")
			if pkg.FileExists(abstractFile) {
				continue
			}
			results = append(results, path.Join(workDir, entry.Name()))
		}
	}
	return results, nil
}

func (s *ScienceDirect) process(path string) error {
	byteData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	jsonResult := gjson.Parse(string(byteData))
	link := jsonResult.Get("link").String()

	targetLink := s.domain + link
	s.log.Infof("Do process link: %s", targetLink)

	err = s.page.Navigate(targetLink)
	if err != nil {
		s.log.Errorf("page.Navigate err: %+v", err)
		return err
	}

	second := time.Duration(time.Now().Unix()%10 + 10)
	time.Sleep(second * time.Second)

	htmlContent, err := s.page.HTML()
	if err != nil {
		return err
	}
	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return err
	}

	content, err := parseText(doc)
	if err != nil {
		return err
	}

	content = strings.TrimLeft(content, "Abstract")

	abstractFile := strings.ReplaceAll(path, ".json", ".abstract")

	return os.WriteFile(abstractFile, []byte(content), 0755)
}

func parseText(doc *html.Node) (string, error) {

	xpaths := []string{
		"//div[@id='abstracts']//text()",
		"//div[@id='body']//text()"}

	for _, xpath := range xpaths {
		contentNode, err := htmlquery.QueryAll(doc, xpath)
		if err != nil {
			continue
		}
		if len(contentNode) == 0 {
			continue
		}

		var article strings.Builder
		for _, node := range contentNode {
			item := htmlquery.InnerText(node)
			article.WriteString(item)
		}
		return article.String(), nil
	}
	return "", nil
}
