package recommendationConst

const (
	AllType      = "all"
	ReceivedType = "received"
	SentType     = "sent"
)

func AllowedRecommendationTypes() []string {
	return []string{
		AllType,
		ReceivedType,
		SentType,
	}
}
