package repositoryModel

type User struct {
	Metadata
	Username       *string `gorm:"uniqueIndex;not null"`
	DisplayName    string
	Email          string `gorm:"uniqueIndex"`
	HashedPassword string
}
