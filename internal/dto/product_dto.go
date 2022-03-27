package dto

type ProductRequest struct {
	Name string `json:"name"`
	Price string `json:"price"`
	Image string `json:"image"`
	Colors string `json:"colors"`
	Compare bool `json:"compare"`
}

type ProductResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Price string `json:"price"`
	Image string `json:"image"`
	Colors string `json:"colors"`
	Compare bool `json:"compare"`
}
