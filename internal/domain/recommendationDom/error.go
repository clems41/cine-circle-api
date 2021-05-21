package recommendationDom

import "cine-circle/internal/typedErrors"

const (
	errSenderIDNullCode     = "SENDER_ID_NULL"
	errMovieIDNullCode      = "MOVIE_ID_NULL"
	errCommentEmptyCode     = "COMMENT_EMPTY"
	errMissingRecipientCode = "MISSING_RECIPIENT"
	errUserUnauthorizedCode = "USER_UNAUTHORIZED"
	errMovieNotFoundCode    = "MOVIE_NOT_FOUND"

	errRecommendationTypeIncorrectCode = "RECOMMENDATION_TYPE_INCORRECT"
	errSortingFieldIncorrectCode       = "SORTING_FIELD_INCORRECT"
)

var (
	errSenderIDNull     = typedErrors.NewBadRequestWithCode(errSenderIDNullCode)
	errMovieIDNull      = typedErrors.NewBadRequestWithCode(errMovieIDNullCode)
	errCommentEmpty     = typedErrors.NewBadRequestWithCode(errCommentEmptyCode)
	errMissingRecipient = typedErrors.NewBadRequestWithCode(errMissingRecipientCode)
	errUserUnauthorized = typedErrors.NewAuthenticationErrorWithCode(errUserUnauthorizedCode)
	errMovieNotFound    = typedErrors.NewNotFoundWithCode(errMovieNotFoundCode)

	errRecommendationTypeIncorrect = typedErrors.NewBadRequestWithCode(errRecommendationTypeIncorrectCode)
	errSortingFieldIncorrect       = typedErrors.NewBadRequestWithCode(errSortingFieldIncorrectCode)
)
