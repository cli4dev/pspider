package main

import (
	"fmt"
	"pspider/spiders"
	"pspider/spiders/tmall"

	"github.com/micro-plat/lib4go/logger"
)

func main() {

	defer logger.Close()
	spider := tmall.NewTmSpider("车载u盘", func(p *spiders.Product) {
		fmt.Println(p)
	})
	if err := spider.Start(); err != nil {
		fmt.Println(err)
	}
}
