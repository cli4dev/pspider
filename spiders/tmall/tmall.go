package tmall

import (
	"fmt"
	"pspider/spiders"
)

//TmSpider 天猫商品爬虫
type TmSpider struct {
	kw         string
	isStart    bool
	notifyChan chan *spiders.Product
	callback   func(*spiders.Product)
}

//NewTmSpider 天猫商品查询
func NewTmSpider(kw string, callback func(*spiders.Product)) *TmSpider {
	if callback == nil {
		panic("tmspider.callback不能为空")
	}
	tm := &TmSpider{
		kw:         kw,
		notifyChan: make(chan *spiders.Product),
		callback:   callback,
	}
	go tm.notify()
	return tm
}
func (t *TmSpider) notify() {
	for {
		select {
		case p, ok := <-t.notifyChan:
			if !ok {
				return
			}
			t.callback(p)
		}
	}
}

//Start 搜索商品信息
func (t *TmSpider) Start() error {
	if t.isStart {
		return nil
	}
	t.isStart = true
	list, err := getProductList(t.kw, 100)
	if err != nil {
		return err
	}
	for _, url := range list[0:1] {
		p, err := getProductDetail("https:" + url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%+v\n", p)
	}
	return nil
}
