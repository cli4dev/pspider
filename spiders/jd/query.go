package jd

import (
	"bytes"
	"context"
	"fmt"
	"pspider/spiders"
	"pspider/spiders/chrome"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/zkfy/log"
)

//getURLs 根据关键字查询商品链接
func (t *jd) getURLs(kw string) error {
	defer close(t.urls)
	err := chrome.Run(chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("https://search.jd.com/Search?keyword=%s&enc=utf-8&wq=%s&pvid=3d9adc305f784919ae876304b675a07d", kw, kw)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 1; i++ {
				t.log.Debugf("查询第%d页数据", i+1)
				chromedp.WaitVisible(`#J_goodsList`).Do(ctx)
				chromedp.ScrollIntoView("#J_bottomPage").Do(ctx)
				chromedp.Sleep(time.Second * 5).Do(ctx)
				if err := t.getURL(ctx); err != nil {
					t.log.Error(err)
				}
				chromedp.Click(".pn-next").Do(ctx)
			}
			return nil
		}),
	}, time.Minute*3, t.log)
	if err != nil {
		t.log.Error(err)
	}
	return err
}
func (t *jd) getURL(ctx context.Context) error {
	var html string
	chromedp.OuterHTML(`#J_goodsList`, &html, chromedp.ByQuery).Do(ctx)
	var b bytes.Buffer
	b.WriteString(html)
	doc, err := goquery.NewDocumentFromReader(&b)
	if err != nil {
		return err
	}
	doc.Find(".gl-item").Each(func(i int, s *goquery.Selection) {
		product := &spiders.Product{}
		if href := chrome.GetHref(s, "div.p-img a"); href != "" {
			product.URL = href
		}
		product.ShopName = chrome.GetTextBySelection(s, ".p-shop a")
		product.ShopURL = chrome.GetHref(s, ".p-shop a")
		product.Appraise = chrome.GetTextBySelection(s, ".p-commit strong a")
		product.MonthlySales = chrome.GetTextBySelection(s, ".p-commit strong a")
		t.urls <- product

	})
	return nil
}

func (t *jd) getProds() error {
	var group sync.WaitGroup
	for i := 0; i < 10; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for {
				select {
				case prod, ok := <-t.urls:
					if !ok {
						return
					}
					p, err := t.getDetails(prod)
					if err != nil {
						log.Error(err)
						continue
					}
					t.notify <- p
				}
			}

		}()
	}
	group.Wait()
	return nil
}

func (t *jd) getDetails(product *spiders.Product) (sp *spiders.Product, err error) {
	err = chrome.Run(chromedp.Tasks{
		chromedp.Navigate(product.URL),
		chromedp.Sleep(time.Second * 5),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var html string
			chromedp.OuterHTML(`.product-intro`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.Title = chrome.GetText(doc, ".sku-name")
			product.SalesPrice = doc.Find(".p-price").Text()
			product.Area = chrome.GetText(doc, ".summary-service")
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var html string
			chromedp.OuterHTML(`.jdm-toolbar-tabs`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.WangWang = chrome.GetHrefByDoc(doc, ".jdm-tbar-tab-contact a")
			return nil
		}),
		chromedp.ScrollIntoView("#footmark"),
		chromedp.WaitVisible(".comment-percent .percent-con"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var html string
			chromedp.OuterHTML(`#comment`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.Score = chrome.GetText(doc, ".comment-percent .percent-con")
			return nil
		}),
	}, time.Minute, t.log)
	if err != nil {
		return nil, err
	}

	return product, nil
}
