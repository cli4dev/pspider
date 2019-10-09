package spiders

import (
	"fmt"

	"github.com/micro-plat/lib4go/logger"
)

type FSpider interface {
	Query(kw string, notify chan *Product) error
}

var spds = map[string]FSpider{}

func Register(name string, f FSpider) {
	if _, ok := spds[name]; ok {
		panic(fmt.Sprintf("%s已注册，请不要重复注册", name))
	}
	spds[name] = f
}

//Spider 商品爬虫
type Spider struct {
	kw         string
	index      int
	log        logger.ILogger
	notifyChan chan *Product
}

//NewSpider 天猫商品查询
func NewSpider(kw string) *Spider {
	sp := &Spider{
		kw:         kw,
		notifyChan: make(chan *Product, 500),
		log:        logger.New("spider"),
	}
	go sp.loopNotify()
	return sp
}
func (t *Spider) loopNotify() {
	for {
		select {
		case p, ok := <-t.notifyChan:
			if !ok {
				return
			}
			t.index++
			t.notify(t.index, p)
		}
	}
}

//Start 搜索商品信息
func (t *Spider) Start() error {
	for _, p := range spds {
		if err := p.Query(t.kw, t.notifyChan); err != nil {
			t.log.Error(err)
		}
	}

	return nil
}
func (t *Spider) notify(i int, p *Product) {
	t.log.Debug(i, p)
}
