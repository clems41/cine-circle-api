package repositoryModel

type Recommendation struct {
	Metadata
	SenderID uint     `gorm:"not null;index:idx_recommendation_sender_id"`
	Sender   *User    `gorm:"association_autoupdate:false;association_autocreate:false"`
	MovieID  uint     `gorm:"not null;index:idx_recommendation_movie_id"`
	Movie    *Movie   `gorm:"association_autoupdate:false;association_autocreate:false"`
	Comment  string   `gorm:"not null;default:null"`
	Circles  []Circle `gorm:"many2many:recommendation_circle;association_autoupdate:false;association_autocreate:false"`
	Users    []User   `gorm:"many2many:recommendation_user;association_autoupdate:false;association_autocreate:false"`
}
