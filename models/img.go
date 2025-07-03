package models

import "gorm.io/gorm"

type Img struct {
	gorm.Model
	Url string `json:"url"` // 图片的访问 URL
}
