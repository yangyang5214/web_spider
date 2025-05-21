package pkg

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	page := rod.New().MustConnect().MustPage()
	go page.EachEvent(func(e *proto.NetworkResponseReceived) (stop bool) {

		t.Log(e)

		return
	})()
	page.MustNavigate("https://baidu.com")

	time.Sleep(10 * time.Second)
}

func TestNanosecond(t *testing.T) {
	t.Log(time.Now().Unix() % 10)
}
