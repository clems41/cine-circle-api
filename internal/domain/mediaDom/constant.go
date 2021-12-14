package mediaDom

import (
	"cine-circle-api/internal/constant/languageConst"
	"cine-circle-api/pkg/utils/httpUtils"
)

var (
	mediaIdParameter = httpUtils.Parameter{
		Name:         "mediaId",
		DefaultValue: "",
		Description:  "Id of media to get (can be tv show or movie)",
		DataType:     "string",
		Required:     true,
	}
	languageParameter = httpUtils.Parameter{
		Name:         "lang",
		DefaultValue: languageConst.FrenchLanguage,
		Description:  "Language to use ro fill media fields (overview, title, etc.)",
		DataType:     "string",
		Required:     false,
	}
	keywordParameter = httpUtils.Parameter{
		Name:         "keyword",
		DefaultValue: "",
		Description:  "Keyword to find media",
		DataType:     "string",
		Required:     true,
	}

	defaultQueryParametersValues = map[string]interface{}{
		languageParameter.Name: languageParameter.DefaultValue,
	}
)
