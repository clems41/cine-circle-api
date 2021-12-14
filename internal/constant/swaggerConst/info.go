package swaggerConst

import "cine-circle-api/pkg/httpServer/swagger"

const (
	title        = "Cine circle API"
	description  = "Share movies with your friends !"
	contactName  = "Teasy"
	contactEmail = "contact@teasy.fr"
	version      = "1.0.0"
)

// Define here all tags that can be used with swagger
const (
	UserTag  = "users"
	OtherTag = "other"
)

var (
	tags = map[swagger.TagName]swagger.TagDescription{
		UserTag:  "Managing own user info",
		OtherTag: "Anything else",
	}
)

var (
	Info = swagger.Info{
		Title:        title,
		Description:  description,
		ContactName:  contactName,
		ContactEmail: contactEmail,
		Version:      version,
		Tags:         tags,
	}
)
