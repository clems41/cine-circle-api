package repositoryModel

type User struct {
	Metadata
	Username       *string `gorm:"uniqueIndex;not null"`
	DisplayName    string
	Email          string `gorm:"uniqueIndex;not null;default:null"`
	HashedPassword string `gorm:"not null;default:null"`
}
