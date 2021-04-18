package authenticationDom

type SignIn struct {
	Username 		string 			`json:"username"`
	Password 		string			`json:"password"`
}

type Result bool
