package buff

type Response struct {
	Code string       `json:"code"`
	Msg  string       `json:"msg"`
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Items      []Item `json:"items"`
	PageNum    int    `json:"page_num"`
	PageSize   int    `json:"page_size"`
	TotalCount int    `json:"total_count"`
	TotalPage  int    `json:"total_page"`
}

type Item struct {
	Appid          int     `json:"appid"`
	Name           string  `json:"name"`
	MarketHashName string  `json:"market_hash_name"`
	ID             int     `json:"id"`
	SellMinPrice   float64 `json:"sell_min_price,string"`
	BuyMaxPrice    float64 `json:"buy_max_price,string"`
}
