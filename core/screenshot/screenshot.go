package screenshot

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/chromedp/chromedp"
	"os"
	"time"
	"weblive/common/mlogger"
	"weblive/common/options"
)

func getMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))[8:24]
}

func apply() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// 禁用GPU，不显示GUI
		chromedp.DisableGPU,
		// 取消沙盒模式
		chromedp.NoSandbox,
		// 隐身模式启动
		chromedp.Flag("incognito", true),
		// 忽略证书错误
		chromedp.Flag("ignore-certificate-errors", true),
		// 窗口最大化
		chromedp.Flag("start-maximized", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		// 禁用扩展
		chromedp.Flag("disable-extensions", true),
		// 禁止加载所有插件
		chromedp.Flag("disable-plugins", true),
		// 禁用浏览器应用
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1024, 768),
	)
	if options.CurrentOption.BrowserPath != "" {
		opts = append(opts, chromedp.ExecPath(options.CurrentOption.BrowserPath))
	}
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(ctx)
	ctx, cancel = context.WithTimeout(ctx, time.Second*20)
	return ctx, cancel
}

func DoScreen(url string, dir string) (string, error) {
	var b2 []byte
	ctx, cancel := apply()
	defer cancel()
	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1024, 768),
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&b2),
	); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir+"/screen/", 0755); err != nil {
		mlogger.Error("Create folder 'screen' fail")
		return "", err
	}
	ws := "screen/" + getMD5Encode(url) + ".png"
	if err := os.WriteFile(dir+"/"+ws, b2, 0755); err != nil {
		mlogger.Error("Write png file fail")
		return "", err
	}
	return ws, nil
}
