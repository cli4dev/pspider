package spiders

type Product struct {
	URL          string
	Title        string
	ShopName     string
	ShopURL      string
	SalesPrice   string
	Area         string
	MonthlySales string
	Appraise     string
	Score        string
	Points       string
}

type ISpider interface {
	Search(kw string, p chan *Product) error
}
