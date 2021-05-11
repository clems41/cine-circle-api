package repositoryModel

type Circle struct {
	Metadata
	Name 			string
	Description 	string
	Users 			[]User `gorm:"many2many:circle_user;"`
}
