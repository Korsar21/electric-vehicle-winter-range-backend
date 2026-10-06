package serializer

type ElectricCarLoadResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	PowerDrawW  int    `json:"power_draw_w"`
	Priority    string `json:"priority"`

	Status    string `json:"status"`
	LikeCount int    `json:"like_count"`

	ImageURL string `json:"image_url"`
	VideoURL string `json:"video_url"`

	IsCreator int `json:"is_creator"`
	IsLiked   int `json:"is_liked"`

	HasNext bool `json:"has_next"`
}

type ElectricCarLoadPublishRequest struct {
	Description string `json:"description"`
	Category    string `json:"category"`
	PowerDrawW  int    `json:"power_draw_w"`
	Priority    string `json:"priority"`
}

type LikeRequest struct {
	Like int `json:"like"`
}
