package repositoryModel

type Watchlist struct {
	Metadata
	UserID  uint   `gorm:"not null;index:idx_watchlist_user_id"`
	User    *User  `gorm:"association_autoupdate:false;association_autocreate:false"`
	MovieID uint   `gorm:"not null;index:idx_watchlist_movie_id"`
	Movie   *Movie `gorm:"association_autoupdate:false;association_autocreate:false"`
}
