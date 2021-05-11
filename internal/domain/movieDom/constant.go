package movieDom

import "cine-circle/internal/utils"

const (
	ExternalApiUrl = "http://www.omdbapi.com/"

	defaultAPIKey = "9d8fa748"
	envAPIKey = "EXTERNAL_API_KEY"

	defaultPlotValue = "full" //(full or short)
	envPlotValue = "EXTERNAL_PLOT_VALUE"

	MovieMedia = "movie"
	SeriesMedia = "series"

	ReleasedLayout = "02 Jan 2006"

	StringArraySeparator = ","
	RunTimeUnit = " min"
)

var (
	ExternalApiKey = utils.GetDefaultOrFromEnv(defaultAPIKey, envAPIKey)
	PlotValue = utils.GetDefaultOrFromEnv(defaultPlotValue, envPlotValue)
)
