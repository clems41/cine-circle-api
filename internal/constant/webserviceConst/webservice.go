package webserviceConst

// Messages used to define error kind
const (
	BadRequestMessage          = "Bad request, some fields are wrong or missing"
	UnauthorizedMessage        = "Unauthorized user to access this route"
	ForbiddenMessage           = "User doesn't have access to this resource"
	UnprocessableEntityMessage = "Body request cannot be proceed"
	NotFoundMessage            = "Resource cannot be found"
)

// Query parameters
const (
	PageQueryName        = "page"
	PageQueryDescription = "Page to return"
	PageQueryType        = "int"

	PageSizeQueryName        = "pageSize"
	PageSizeQueryDescription = "Number of element by page"
	PageSizeQueryType        = "int"

	SortQueryName        = "sort"
	SortQueryDescription = "Sort result array using specific fields (more than one can be used)"
	SortQueryType        = "string"
)
