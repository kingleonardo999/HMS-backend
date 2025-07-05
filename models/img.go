package models

type Img struct {
	ID  uint   `json:"id" gorm:"primary_key"`
	Url string `json:"url"` // 图片的访问 URL
}
