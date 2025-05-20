package pkg

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"io"
	"net/http"
	"os"
	"time"
)

type ChromePool struct {
	Browser *rod.Browser
	log     *log.Helper
}

func NewChromePool(logger log.Logger, ws bool) (*ChromePool, func(), error) {
	llog := log.NewHelper(logger)

	var (
		launcherURL string
		err         error
	)

	if ws {
		llog.Infof("ws is true, use ws url")
		launcherURL, err = parseUrl()
	} else {
		useDataDir := "/tmp/chrome"

		chromeLauncher := launcher.New().
			Leakless(true).
			Headless(false).
			Env(append(os.Environ(), "TZ=Asia/Shanghai")...).
			UserDataDir(useDataDir)

		launcherURL, err = chromeLauncher.Launch()
	}

	if err != nil {
		llog.Errorf("error: %v", err)
		return nil, nil, err
	}

	llog.Infof("Use chrome url: %s", launcherURL)

	browser := rod.New().ControlURL(launcherURL)
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
