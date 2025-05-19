package sciencedirect

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/tidwall/gjson"
	"net/url"
	"os"
	"path"
	"strings"
	"web_spider/pkg"
)

type ScienceDirect struct {
	log *log.Helper

	chrome  *pkg.ChromePool
	workDir string
	page    *rod.Page
}

func NewScienceDirect(chrome *pkg.ChromePool, logger log.Logger) *ScienceDirect {
	homeDir, _ := os.UserHomeDir()

	return &ScienceDirect{
		log:     log.NewHelper(logger),
		chrome:  chrome,
		workDir: path.Join(homeDir, "sciencedirect"),
	}
}

func (s *ScienceDirect) List() error {
	page, err := s.chrome.Browser.Page(proto.TargetCreateTarget{
		URL: "https://www.sciencedirect.com/",
	})

	if err != nil {
		return err
	}
	s.page = page

	go page.EachEvent(func(e *proto.NetworkResponseReceived) {
		//save to dir
		respUrl := e.Response.URL
		if strings.HasPrefix(respUrl, "https://www.sciencedirect.com/search/api") {
			s.log.Infof("sciencedirect api url: %s", respUrl)
			go s.saveDir(respUrl, e.RequestID)
		}
	})()

	select {}
}

func (s *ScienceDirect) saveDir(targetUrl string, reqId proto.NetworkRequestID) (err error) {
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

		_, err = os.Stat(resultPath)
		if errors.Is(err, os.ErrNotExist) {
			continue // skip if file exists
		}

		err = os.WriteFile(resultPath, []byte(sr.String()), 0755)
		if err != nil {
			s.log.Errorf("save to local dir err: %+v", err)
			return err
		}
	}
	return
}

func (s *ScienceDirect) Detail() error {
	return nil
}
