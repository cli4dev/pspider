package main

import (
	"fmt"
	"pspider/spiders"

	_ "pspider/spiders/tmall"

	_ "pspider/spiders/jd"

	"github.com/micro-plat/lib4go/logger"
)

func main() {

	defer logger.Close()
	spider := spiders.NewSpider("充值系统")
	if err := spider.Start(); err != nil {
		fmt.Println(err)
	}
}
