package recommendationDom

import (
	typedErrors2 "cine-circle/pkg/typedErrors"
)

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
	errSenderIDNull     = typedErrors2.NewBadRequestWithCode(errSenderIDNullCode)
	errMovieIDNull      = typedErrors2.NewBadRequestWithCode(errMovieIDNullCode)
	errCommentEmpty     = typedErrors2.NewBadRequestWithCode(errCommentEmptyCode)
	errMissingRecipient = typedErrors2.NewBadRequestWithCode(errMissingRecipientCode)
	errUserUnauthorized = typedErrors2.NewAuthenticationErrorWithCode(errUserUnauthorizedCode)
	errMovieNotFound    = typedErrors2.NewNotFoundWithCode(errMovieNotFoundCode)

	errRecommendationTypeIncorrect = typedErrors2.NewBadRequestWithCode(errRecommendationTypeIncorrectCode)
	errSortingFieldIncorrect       = typedErrors2.NewBadRequestWithCode(errSortingFieldIncorrectCode)
)
