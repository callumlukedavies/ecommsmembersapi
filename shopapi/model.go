package shopapi

type Item struct {
	ID           string `json:"ItemID"`
	Name         string `json:"ItemName"`
	Gender       string `json:"ItemGender"`
	Description  string `json:"ItemDescription"`
	ImageName    string `json:"ItemImageName"`
	GalleryImage string `json:"ItemGalleryImage"`
	DateUploaded string `json:"ItemUploadDate"`
	Price        string `json:"ItemPrice"`
	IsSold       bool   `json:"ItemIsSold"`
	Size         string `json:"ItemSize"`
	Category     string `json:"ItemCategory"`
	Condition    string `json:"ItemCondition"`
	SellerID     int    `json:"ItemSellerID"`
	SellerName   string `json:"ItemSellerName"`
}
