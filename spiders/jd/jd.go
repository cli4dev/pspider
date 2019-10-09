package jd

import (
	"pspider/spiders"
	"time"

	"github.com/micro-plat/lib4go/logger"
)

type jd struct {
	urls   chan *spiders.Product
	log    logger.ILogger
	notify chan *spiders.Product
}

func newjd() *jd {
	return &jd{
		log:  logger.New("jd"),
		urls: make(chan *spiders.Product, 10000),
	}
}

func (t *jd) Query(kw string, notify chan *spiders.Product) error {
	t.notify = notify
	start := time.Now()

	t.log.Infof("1. 从京东抓取商品信息[%s]", kw)
	go t.getURLs(kw)
	err := t.getProds()

	t.log.Info("2. 京东商品抓取完成", time.Since(start))
	return err
}

func init() {
	spiders.Register("jd", newjd())
}
