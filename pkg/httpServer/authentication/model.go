package authentication

// UserInfo specify which information are needed for the authentication (don't add unnecessary fields).
// !!! We don't want to use model defined as DTO (in domain) or model defined as EO (in repository) !!!
type UserInfo struct {
	Id             uint
	Role           string
}
