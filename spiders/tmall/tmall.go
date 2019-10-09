package tmall

import (
	"pspider/spiders"
	"time"

	"github.com/micro-plat/lib4go/logger"
)

type tmall struct {
	urls   chan *spiders.Product
	log    logger.ILogger
	notify chan *spiders.Product
}

func newTmall() *tmall {
	return &tmall{
		log:  logger.New("tmall"),
		urls: make(chan *spiders.Product, 10000),
	}
}

func (t *tmall) Query(kw string, notify chan *spiders.Product) error {
	t.notify = notify
	start := time.Now()

	t.log.Infof("1. 从天猫抓取商品信息[%s]", kw)
	go t.getURLs(kw)
	err := t.getProds()

	t.log.Info("2. 天猫商品抓取完成", time.Since(start))
	return err
}

func init() {
	spiders.Register("tmall", newTmall())
}
