package recommendationDom

const (
	recommendationReceived = "received"
	recommendationSent     = "sent"
	recommendationBoth     = ""
)

var (
	acceptedTypeOfRecommendation = []string{recommendationReceived, recommendationSent, recommendationBoth}
	acceptedFieldsForSorting     = []string{"date", "recommendationType"}
)
