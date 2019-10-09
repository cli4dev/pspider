package chrome

import (
	"context"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/micro-plat/lib4go/logger"
)

func Run(task chromedp.Tasks, t time.Duration, log logger.ILogger) error {
	opts := make([]chromedp.ExecAllocatorOption, 0)

	opts = append(opts, chromedp.Flag("headless", true))
	opts = append(opts, chromedp.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"))
	ctx, timeout := context.WithTimeout(context.Background(), t)
	defer timeout()
	allocatorCtx, _ := chromedp.NewExecAllocator(ctx, opts...)
	ctx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Errorf))
	defer cancel()
	err := chromedp.Run(ctx, task)
	if err != nil {
		return err
	}
	return nil
}
func GetText(doc *goquery.Document, path string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.Replace(doc.Find(path).Text(), "\n", "", -1)), " ", "")
}

func GetTextBySelection(a *goquery.Selection, path string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.Replace(a.Find(path).Text(), "\n", "", -1)), " ", "")
}
func GetAttr(doc *goquery.Document, path string, attr string) string {
	v, _ := doc.Find(path).Attr(attr)
	return strings.TrimSpace(strings.Replace(v, "\n", "", -1))
}
func GetHrefByDoc(doc *goquery.Document, path string) string {
	href, ok := doc.Find(path).Attr("href")
	if ok && !strings.HasPrefix(href, "http") {
		return "https:" + href
	}
	return href
}
func GetHref(s *goquery.Selection, path string) string {
	a := s.Find(path)
	href, ok := a.Attr("href")
	if ok && !strings.HasPrefix(href, "http") {
		return "https:" + href
	}
	return href
}

func Wait(path string, ctx context.Context, timeout time.Duration) {
	for {
		select {
		case <-time.After(timeout):
			return
		default:
			value := ""
			chromedp.Text(path, &value, chromedp.ByID).Do(ctx)
			if value == "" {
				time.Sleep(time.Second)
			}
			return
		}
	}
}
