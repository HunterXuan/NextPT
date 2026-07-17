package economyin

type ShopProductListInp struct{}

type ShopOrderCreateInp struct {
	ProductKey string `json:"productKey" v:"required|max-length:64" description:"商品标识"`
}

type ShopOrderListInp struct {
	Page int `json:"page" d:"1" v:"min:1" description:"页码"`
	Size int `json:"size" d:"20" v:"min:1|max:100" description:"每页数量"`
}
