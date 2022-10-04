package db

import "time"

type Url struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	ShortId   string `json:"short_id" gorm:"short_id"`
	Url       string `json:"long_url" gorm:"url"`
	Clicks    uint   `json:"clicks" gorm:"clicks,default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
