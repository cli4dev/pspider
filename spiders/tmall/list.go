package tmall

import (
	"bytes"
	"context"
	"fmt"
	"pspider/spiders"
	"pspider/spiders/chrome"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
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
				if href, ok := s.Attr("href"); ok {
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

func getProductDetail(url string) (*spiders.Product, error) {

	product := &spiders.Product{}
	err := chrome.Run(chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#detail .tm-count`),
		chromedp.Sleep(time.Second * 10),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var html string
			chromedp.OuterHTML(`#detail`, &html, chromedp.ByQuery).Do(ctx)
			doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
			if err != nil {
				return err
			}
			product.Title = chrome.GetText(doc, ".tb-detail-hd h1")
			product.SalesPrice = doc.Find(".tm-price").Text()
			product.MonthlySales = chrome.GetText(doc, ".module-wrap .sales")
			product.Appraise = chrome.GetText(doc, "#J_ItemRates .tm-indcon")
			product.Score = chrome.GetText(doc, ".tm-ind-emPointCount .tm-count")
			return nil
		}),
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}
