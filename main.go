package main

import (
	"fmt"
	"pspider/spiders"
	"pspider/spiders/tmall"
)

func main() {

	spider := tmall.NewTmSpider("车载u盘", func(p *spiders.Product) {
		fmt.Println(p)
	})
	fmt.Println(spider.Start())
}
