package mediaDom

import (
	"cine-circle-api/internal/constant/mediaConst"
	"cine-circle-api/pkg/utils/httpUtils"
)

var (
	mediaIdParameter = httpUtils.Parameter{
		Name:            "mediaId",
		DefaultValueStr: "",
		Description:     "Id of media to get (can be tv show or movie)",
		DataType:        "string",
		Required:        true,
	}
	languageParameter = httpUtils.Parameter{
		Name:            "lang",
		DefaultValueStr: mediaConst.FrenchLanguage,
		Description:     "Language to use ro fill media fields (overview, title, etc.)",
		DataType:        "string",
		Required:        false,
	}
	keywordParameter = httpUtils.Parameter{
		Name:            "keyword",
		DefaultValueStr: "",
		Description:     "Keyword to find media",
		DataType:        "string",
		Required:        true,
	}

	defaultQueryParametersValues = map[string]string{
		languageParameter.Name: languageParameter.DefaultValueStr,
	}
)
