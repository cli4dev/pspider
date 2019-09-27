package spiders

type Product struct {
	Title         string
	Payment       string
	ShopName      string
	OriginalPrice string
	SalesPrice    string
	Area          string
	PostFee       string
	MonthlySales  string
	Appraise      string
	Score         string
}

type ISpider interface {
	Search(kw string, p chan *Product) error
}
