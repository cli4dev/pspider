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
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
)

//getProductList 根据关键字查询商品列表，返回商品链接数组
//1. 根据关健字查询
//2. 点击排序
//3. 检查页面
//4. 根据页面进行链接获取
//5. 查询所有页码数据后返回结果
func getProductList(kw string, count int, orderBy ...string) ([]string, error) {

	list := make([]string, 0, 100)
	err := chrome.Run(chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("https://list.tmall.com/search_product.htm?q=%s&from=mallfp..pc_1_searchbutton", kw)),
		chromedp.WaitVisible(`#J_ItemList`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var html string
			chromedp.OuterHTML(`#J_ItemList`, &html, chromedp.ByQuery).Do(ctx)
			var b bytes.Buffer
			b.WriteString(html)
			doc, err := goquery.NewDocumentFromReader(&b)
			if err != nil {
				return err
			}
			doc.Find("div.productImg-wrap a").Each(func(i int, s *goquery.Selection) {
				if href := chrome.GetHref(s); href != "" {
					list = append(list, href)
				}
			})
			return nil
		}),
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func getProducts(log logger.ILogger, urls ...string) ([]*spiders.Product, error) {
	ps := make([]*spiders.Product, 0, len(urls))
	var group sync.WaitGroup
	ch := make(chan string, len(urls))
	min := types.GetMin(len(urls), 10)
	for _, url := range urls {
		ch <- url
	}
	count := 0
	for i := 0; i < min; i++ {
		go func() {
			for {
				url, ok := <-ch
				if !ok {
					return
				}
				group.Add(1)
				p, err := getProductDetail(url)
				if err != nil {
					count++
					group.Done()
					log.Error(err)
					continue
				}
				count++
				log.Info("完成:", count, url)
				ps = append(ps, p)
				group.Done()
			}

		}()
	}
	time.Sleep(time.Second)

	group.Wait()
	close(ch)
	return ps, nil
}

func getProductDetail(url string) (*spiders.Product, error) {
	product := &spiders.Product{}
	err := chrome.Run(chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(time.Second * 15),
		chromedp.ActionFunc(func(ctx context.Context) error {

			var html string
			chromedp.OuterHTML(`#header`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.ShopName = chrome.GetText(doc, ".slogo-shopname")
			product.ShopURL = chrome.GetAttr(doc, ".slogo-shopname", "href")
			product.Score = chrome.GetText(doc, "#shop-info .main-info")
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var html string
			chromedp.OuterHTML(`#detail`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.URL = url
			product.Title = chrome.GetText(doc, ".tb-detail-hd h1")
			product.SalesPrice = doc.Find(".tm-price").Text()
			product.MonthlySales = chrome.GetText(doc, ".tm-ind-sellCount .tm-count")
			product.Appraise = chrome.GetText(doc, ".tm-ind-reviewCount .tm-count")
			product.Points = chrome.GetText(doc, ".tm-ind-emPointCount .tm-count")
			return nil
		}),
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}
