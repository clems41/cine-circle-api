package recommendationDom

import (
	"cine-circle-api/internal/constant/recommendationConst"
	"cine-circle-api/pkg/utils/httpUtils"
	"fmt"
)

// Path or query parameters
var (
	pageQueryParameter = httpUtils.Parameter{
		Name:            "page",
		Description:     "Page number to return",
		DefaultValueStr: "1",
		DataType:        "int",
		Required:        false,
	}
	pageSizeQueryParameter = httpUtils.Parameter{
		Name:            "pageSize",
		Description:     "Number of element by page",
		DefaultValueStr: "10",
		DataType:        "int",
		Required:        false,
	}
	recommendationTypeQueryParameter = httpUtils.Parameter{
		Name:            "type",
		Description:     fmt.Sprintf("Filter on type of recommendation %v", recommendationConst.AllowedRecommendationTypes()),
		DefaultValueStr: "all",
		DataType:        "string",
		Required:        false,
	}
	movieIdQueryParameter = httpUtils.Parameter{
		Name:            "movieId",
		Description:     "Filter on specific movie",
		DefaultValueStr: "",
		DataType:        "int",
		Required:        false,
	}
	defaultSearchQueryParameterValues = map[string]string{
		pageQueryParameter.Name:               pageQueryParameter.DefaultValueStr,
		pageSizeQueryParameter.Name:           pageSizeQueryParameter.DefaultValueStr,
		recommendationTypeQueryParameter.Name: recommendationTypeQueryParameter.DefaultValueStr,
		movieIdQueryParameter.Name:            movieIdQueryParameter.DefaultValueStr,
	}
)
