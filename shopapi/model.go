package shopapi

type Item struct {
	Name      string `json:"ItemName"`
	ImageName string `json:"ItemImageName"`
	ID        int    `json:"ItemID"`
	Price     int    `json:"ItemPrice"`
	SellerID  int    `json:"ItemSellerID"`
}
