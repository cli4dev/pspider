package tmall

import (
	"pspider/spiders"
	"time"

	"github.com/micro-plat/lib4go/logger"
)

//TmSpider 天猫商品爬虫
type TmSpider struct {
	kw         string
	isStart    bool
	log        logger.ILogger
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
		log:        logger.New("tmspider"),
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
	start := time.Now()
	t.log.Info("-----------开始抓取天猫商品数据-----------")
	t.log.Infof("1. 查询商品列表[%s %d]", t.kw, 100)
	list, err := getProductList(t.kw, 100)
	if err != nil {
		return err
	}
	t.log.Infof("2. 查询到%d个商品，开始抓取商品明细", len(list))
	ps, err := getProducts(list...)
	if err != nil {
		return err
	}
	t.log.Infof("3. 抓取完成,总共获取到%d个商品 %v", len(ps), time.Since(start))
	return nil
}
