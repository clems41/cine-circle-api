package circleDom

import (
	"cine-circle-api/pkg/utils/httpUtils"
)

// Path or query parameters
var (
	circleIdPathParameter = httpUtils.Parameter{
		Name:            "circleId",
		Description:     "Id of the resource",
		DefaultValueStr: "",
		DataType:        "int",
		Required:        true,
	}
	userIdPathParameter = httpUtils.Parameter{
		Name:            "userId",
		Description:     "Id of the user to add/delete from circle",
		DefaultValueStr: "",
		DataType:        "int",
		Required:        true,
	}
	circleNameQueryParameter = httpUtils.Parameter{
		Name:            "name",
		Description:     "Name of the circle to find",
		DefaultValueStr: "",
		DataType:        "string",
		Required:        true,
	}
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
	defaultSearchQueryParameterValues = map[string]string{
		pageQueryParameter.Name:     pageQueryParameter.DefaultValueStr,
		pageSizeQueryParameter.Name: pageSizeQueryParameter.DefaultValueStr,
	}
)
