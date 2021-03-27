package model

type Watchlist struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
	MovieID string `json:"movie_id" gorm:"primaryKey"`
}
