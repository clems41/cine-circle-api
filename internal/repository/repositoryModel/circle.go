package repositoryModel

type Circle struct {
	Metadata
	Name        string `gorm:"not null;default:null"`
	Description string `gorm:"not null;default:null"`
	Users       []User `gorm:"many2many:circle_user;association_autoupdate:false;association_autocreate:false"`
}
