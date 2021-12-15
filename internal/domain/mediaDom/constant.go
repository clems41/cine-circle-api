package mediaDom

import (
	"cine-circle-api/pkg/utils/httpUtils"
)

var (
	mediaIdPathParameter = httpUtils.Parameter{
		Name:            "mediaId",
		DefaultValueStr: "",
		Description:     "Id of media to get",
		DataType:        "int",
		Required:        true,
	}
	keywordQueryParameter = httpUtils.Parameter{
		Name:            "keyword",
		DefaultValueStr: "",
		Description:     "Keyword to find media",
		DataType:        "string",
		Required:        true,
	}
	pageQueryParameter = httpUtils.Parameter{
		Name:            "page",
		Description:     "Current page to display",
		DefaultValueStr: "1",
		DataType:        "int",
		Required:        false,
	}
	pageSizeQueryParameter = httpUtils.Parameter{
		Name:            "pageSize",
		Description:     "Number of element by page",
		DefaultValueStr: "20",
		DataType:        "int",
		Required:        false,
	}
	defaultSearchQueryParameterValues = map[string]string{
		pageQueryParameter.Name:     pageQueryParameter.DefaultValueStr,
		pageSizeQueryParameter.Name: pageSizeQueryParameter.DefaultValueStr,
	}
)
