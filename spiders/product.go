package spiders

import "encoding/json"

//Product 商品信息
type Product struct {
	//商品地址
	URL string `json:"url"`

	//商品标题
	Title string `json:"title"`

	//店铺名称
	ShopName string `json:"shop-name"`

	//店铺地址
	ShopURL string `json:"shop-url"`

	//售价
	SalesPrice string `json:"sale-price"`

	//地区
	Area string `json:"area"`

	//月销售量
	MonthlySales string `json:"monthly-sales"`

	//评价总数
	Appraise string `json:"appraise"`

	//商品得分
	Score string `json:"score"`

	//联系方式
	WangWang string `json:"ww"`
}

type ISpider interface {
	Search(kw string, p chan *Product) error
}

func (p *Product) String() string {
	text, _ := json.Marshal(p)
	return string(text)
}
