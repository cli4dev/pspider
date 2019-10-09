package tmall

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
)

//getURLs 根据关键字查询商品链接
func (t *tmall) getURLs(kw string) error {
	defer close(t.urls)
	err := chrome.Run(chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("https://list.tmall.com/search_product.htm?q=%s&from=mallfp..pc_1_searchbutton", kw)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 1; i++ {
				t.log.Debugf("查询第%d页数据", i+1)
				chromedp.WaitVisible(`#J_ItemList`).Do(ctx)
				chromedp.ScrollIntoView(".ui-page").Do(ctx)
				chromedp.Sleep(time.Second * 5).Do(ctx)
				if err := t.getURL(ctx); err != nil {
					t.log.Error(err)
				}
				chromedp.Click(".ui-page-next").Do(ctx)
			}
			return nil
		}),
	}, time.Minute*3, t.log)
	if err != nil {
		t.log.Error(err)
	}
	return err
}
func (t *tmall) getURL(ctx context.Context) error {
	var html string
	chromedp.OuterHTML(`#J_ItemList`, &html, chromedp.ByQuery).Do(ctx)
	var b bytes.Buffer
	b.WriteString(html)
	doc, err := goquery.NewDocumentFromReader(&b)
	if err != nil {
		return err
	}
	doc.Find("div.product-iWrap").Each(func(i int, s *goquery.Selection) {
		product := &spiders.Product{}
		if href := chrome.GetHref(s, "div.productImg-wrap a"); href != "" {
			product.URL = href
		}
		product.MonthlySales = s.Find("p.productStatus span em").Text()
		t.urls <- product

	})
	return nil
}
func (t *tmall) getProds() error {
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
						t.log.Error(err)
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

func (t *tmall) getDetails(product *spiders.Product) (sp *spiders.Product, err error) {
	err = chrome.Run(chromedp.Tasks{
		chromedp.Navigate(product.URL),
		chromedp.Sleep(time.Second * 15),
		chromedp.ActionFunc(func(ctx context.Context) error {

			var html string
			chromedp.OuterHTML(`#header`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.ShopName = chrome.GetText(doc, ".slogo-shopname")
			product.ShopURL = chrome.GetHrefByDoc(doc, ".slogo-shopname") //chrome.GetAttr(doc, ".slogo-shopname", "href")
			product.Score = chrome.GetText(doc, "#shop-info .main-info")
			product.WangWang = chrome.GetHrefByDoc(doc, "div.slogo-extraicon .ww-static a")
			chromedp.Value("input[name=region]", &product.Area).Do(ctx)

			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var html string
			chromedp.OuterHTML(`#detail`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.Title = chrome.GetText(doc, ".tb-detail-hd h1")
			product.SalesPrice = doc.Find(".tm-price").Text()
			product.Appraise = chrome.GetText(doc, ".tm-ind-reviewCount .tm-count")
			return nil
		}),
	}, time.Minute, t.log)
	if err != nil {
		return nil, err
	}

	return product, nil
}
