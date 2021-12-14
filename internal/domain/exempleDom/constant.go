package exempleDom

import "cine-circle-api/internal/constant/webserviceConst"

// Valeurs par défaut pour le tri des exemples
const (
	defaultPage     = 1
	defaultPageSize = 20
)

var (
	defaultSort                        = []string{"id:asc"}
	defaultSearchQueryParametersValues = map[string]interface{}{
		webserviceConst.PageQueryName:     defaultPage,
		webserviceConst.PageSizeQueryName: defaultPageSize,
		webserviceConst.SortQueryName:     defaultSort,
	}
)
