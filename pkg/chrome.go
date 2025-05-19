package pkg

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod"
	"io"
	"net/http"
	"time"
)

type ChromePool struct {
	Browser *rod.Browser
	log     *log.Helper
}

func NewChromePool(logger log.Logger) (*ChromePool, func(), error) {
	llog := log.NewHelper(logger)

	urlStr, err := parseUrl()
	if err != nil {
		return nil, nil, err
	}

	llog.Infof("Use chrome url: %s", urlStr)

	browser := rod.New().ControlURL(urlStr).Timeout(30 * time.Second)
	err = browser.Connect()
	if err != nil {
		return nil, nil, err
	}

	cancel := func() {
		_ = browser.Close()
	}

	return &ChromePool{
		log:     llog,
		Browser: browser,
	}, cancel, nil
}

func parseUrl() (string, error) {
	httpClient := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := httpClient.Get("http://127.0.0.1:9222/json/version")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(respBytes, &data)
	if err != nil {
		return "", err
	}
	webSocketDebuggerUrl, _ := data["webSocketDebuggerUrl"].(string)
	return webSocketDebuggerUrl, nil
}
