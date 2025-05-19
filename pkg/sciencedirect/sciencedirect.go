package sciencedirect

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod/lib/proto"
	"web_spider/pkg"
)

type ScienceDirect struct {
	log *log.Helper

	chrome *pkg.ChromePool
}

func NewScienceDirect(chrome *pkg.ChromePool, logger log.Logger) *ScienceDirect {
	return &ScienceDirect{
		log:    log.NewHelper(logger),
		chrome: chrome,
	}
}

func (s *ScienceDirect) List() error {
	page := s.chrome.Browser.MustPage()

	page = page.MustNavigate("https://www.sciencedirect.com/")

	go page.EachEvent(func(e *proto.NetworkResponseReceived) {
		s.log.Info(e.RequestID)
	})

	select {}
	return nil
}

func (s *ScienceDirect) Detail() error {
	return nil
}
