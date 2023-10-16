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
	MediaTag = "medias"
	OtherTag = "other"
	CircleTag = "circles"
	RecommendationTag = "recommendations"
)

var (
	tags = map[swagger.TagName]swagger.TagDescription{
		UserTag:  "Managing own user info",
		MediaTag: "Search among medias (movie and tv shows)",
		OtherTag: "Anything else",
		CircleTag: "Operations about circles",
		RecommendationTag: "Send and see recommendations",
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
