package model

import (
	"cine-circle-api/pkg/sql/gormUtils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	gormUtils.Metadata
	Username       string `gorm:"uniqueIndex;check:username <> ''"` // Doit être unique et non vide, car utilisé pour l'authentification
	HashedPassword string `gorm:"check:hashed_password <> ''"`      // Doit être non vide, car utilisé pour l'authentification
	LastName       string `gorm:"check:last_name <> ''"`            // Doit être non vide, car nécessaire pour identifier les utilisateurs
	FirstName      string `gorm:"check:first_name <> ''"`           // Doit être non vide, car nécessaire pour identifier les utilisateurs
	Email          string `gorm:"uniqueIndex;check:email <> ''"`    // Doit être unique et non nul, car utilisé pour l'authentification
	Role           string // Exemple de contrainte : `gorm:"check:role in ('admin', 'lecteur')"` // Doit être égal à un rôle existant : admin ou lecteur
	Active         bool
	EmailToken     string // Utilisé pour valider l'email
	PasswordToken  string // Utilisé pour réinitialiser le mot de passe (mot de passe oublié)
	EmailConfirmed bool   // Utilisé pour valider l'email
}

func MigrateUser(DB *gorm.DB) (err error) {
	err = DB.
		AutoMigrate(&User{})
	if err != nil {
		return errors.WithStack(err)
	}

	return
}
