package com.teasy.CineCircleApi.models.exceptions;

import lombok.Getter;
import org.springframework.http.HttpStatus;

@Getter
public enum ErrorDetails {
    /* Users errors */
    ERR_USER_NOT_FOUND("User %s cannot be found", ErrorOnObject.USER, ErrorOnField.USER_ID, HttpStatus.NOT_FOUND),
    ERR_USER_EMAIL_ALREADY_EXISTS("Email %s already exists and cannot be used twice",
            ErrorOnObject.USER, ErrorOnField.USER_EMAIL, HttpStatus.BAD_REQUEST),
    ERR_USER_USERNAME_ALREADY_EXISTS("Username %s already exists and cannot be used twice",
            ErrorOnObject.USER, ErrorOnField.USER_USERNAME, HttpStatus.BAD_REQUEST),
    ERR_USER_PASSWORD_NOT_MATCHING("Password provided for user with id %s does not match the one in database",
            ErrorOnObject.USER, ErrorOnField.USER_PASSWORD, HttpStatus.BAD_REQUEST),
    ERR_USER_PASSWORD_TOO_SHORT("Password provided %s for user is too short, should be at least 6 characters",
            ErrorOnObject.USER, ErrorOnField.USER_PASSWORD, HttpStatus.BAD_REQUEST),
    ERR_USER_USERNAME_TOO_SHORT("Username provided %s for user is too short, should be at least 4 characters",
            ErrorOnObject.USER, ErrorOnField.USER_USERNAME, HttpStatus.BAD_REQUEST),
    ERR_USER_EMAIL_INCORRECT("Email provided %s does not seem to be a correct email address",
            ErrorOnObject.USER, ErrorOnField.USER_EMAIL, HttpStatus.BAD_REQUEST),
    ERR_USER_RESET_PASSWORD_TOKEN_EMPTY("Reset password token provided is empty",
            ErrorOnObject.USER, ErrorOnField.USER_RESET_PASSWORD_TOKEN, HttpStatus.BAD_REQUEST),
    ERR_USER_RESET_PASSWORD_TOKEN_INCORRECT("Reset password token provided for user with id %s does not match the one in database",
            ErrorOnObject.USER, ErrorOnField.USER_RESET_PASSWORD_TOKEN, HttpStatus.BAD_REQUEST),
    ERR_USER_DISPLAY_NAME_EMPTY("DisplayName cannot be empty for users",
            ErrorOnObject.USER, ErrorOnField.USER_DISPLAY_NAME, HttpStatus.BAD_REQUEST),
    ERR_USER_DISPLAY_NAME_INCORRECT("DisplayName provided %s is incorrect, should be between 4 and 20 characters",
            ErrorOnObject.USER, ErrorOnField.USER_DISPLAY_NAME, HttpStatus.BAD_REQUEST),
    ERR_USER_JWT_TOKEN_EMPTY("Jwt Token provided is empty",
            ErrorOnObject.USER, ErrorOnField.USER_JWT_TOKEN, HttpStatus.BAD_REQUEST),
    ERR_USER_REFRESH_TOKEN_EMPTY("Refresh Token provided is empty",
            ErrorOnObject.USER, ErrorOnField.USER_REFRESH_TOKEN, HttpStatus.BAD_REQUEST),

    /* Library errors */
    ERR_LIBRARY_NOT_FOUND("Library for username %s and mediaId %S cannot be found",
            ErrorOnObject.LIBRARY, null, HttpStatus.NOT_FOUND),
    ERR_LIBRARY_RATING_INCORRECT("Rating provided %s is not correct, should be between 1 and 5",
            ErrorOnObject.LIBRARY, ErrorOnField.LIBRARY_RATING, HttpStatus.BAD_REQUEST),

    /* Heading errors */

    /* Contact errors */
    ERR_CONTACT_FEEDBACK_EMPTY("Feedback provided %s cannot be empty", ErrorOnObject.CONTACT, ErrorOnField.CONTACT_FEEDBACK, HttpStatus.BAD_REQUEST),

    /* Notification errors */

    /* Recommendation errors */
    ERR_RECOMMENDATION_NOT_FOUND("Recommendation with id %s cannot be found",
            ErrorOnObject.RECOMMENDATION, ErrorOnField.RECOMMENDATION_ID, HttpStatus.NOT_FOUND),
    ERR_RECOMMENDATION_TYPE_NOT_SUPPORTED("Recommendation type %s is not supported, should be sent or received",
            ErrorOnObject.RECOMMENDATION, ErrorOnField.RECOMMENDATION_TYPE, HttpStatus.BAD_REQUEST),
    ERR_RECOMMENDATION_RECEIVER_BAD_PERMISSIONS("Authenticated user %s cannot take action on this recommendation id %s",
            ErrorOnObject.RECOMMENDATION, ErrorOnField.RECOMMENDATION_TYPE, HttpStatus.FORBIDDEN),
    ERR_RECOMMENDATION_RATING_INCORRECT("Rating provided %s is not correct, should be between 1 and 5",
            ErrorOnObject.RECOMMENDATION, ErrorOnField.RECOMMENDATION_RATING, HttpStatus.BAD_REQUEST),
    ERR_RECOMMENDATION_USER_IDS_INCORRECT("UserIds should not be empty, you should specify at least one user id",
            ErrorOnObject.RECOMMENDATION, ErrorOnField.RECOMMENDATION_USER_IDS, HttpStatus.BAD_REQUEST),
    ERR_RECOMMENDATION_MEDIA_ID_INCORRECT("mediaId provided %s should be a valid UUID",
            ErrorOnObject.RECOMMENDATION, ErrorOnField.RECOMMENDATION_MEDIA_ID, HttpStatus.BAD_REQUEST),

    /* Token errors */

    /* Watchlist errors */

    /* Media errors */
    ERR_MEDIA_NOT_FOUND("Media with id %s cannot be found", ErrorOnObject.MEDIA, ErrorOnField.MEDIA_ID, HttpStatus.NOT_FOUND),
    ERR_MEDIA_TYPE_NOT_SUPPORTED("Media type %d is not supported", ErrorOnObject.MEDIA, ErrorOnField.MEDIA_MEDIA_TYPE, HttpStatus.BAD_REQUEST),
    ERR_MEDIA_DTO_UNKNOWN("MediaDto class %s is unknown, should never happen", ErrorOnObject.MEDIA, ErrorOnField.MEDIA_ID, HttpStatus.INTERNAL_SERVER_ERROR),

    /* Circle errors */
    ERR_CIRCLE_NOT_FOUND("Circle with id %s cannot be found", ErrorOnObject.CIRCLE, ErrorOnField.CIRCLE_ID, HttpStatus.NOT_FOUND),
    ERR_CIRCLE_USER_BAD_PERMISSIONS("The authenticated user %s has no permissions to update or read this circle",
            ErrorOnObject.CIRCLE, ErrorOnField.CIRCLE_CREATED_BY, HttpStatus.FORBIDDEN),
    ERR_CIRCLE_NAME_EMPTY("Name provided %s cannot be empty", ErrorOnObject.CIRCLE, ErrorOnField.CIRCLE_NAME, HttpStatus.BAD_REQUEST),

    /* Email errors */
    ERR_EMAIL_BUILDING_REQUEST(
            "An error occurred while building the email with id %s from request",
            ErrorOnObject.EMAIL_SERVICE,
            null,
            HttpStatus.INTERNAL_SERVER_ERROR
    ),
    ERR_EMAIL_SENDING_REQUEST(
            "An error occurred while sending the email with id %s to %s",
            ErrorOnObject.EMAIL_SERVICE,
            null,
            HttpStatus.INTERNAL_SERVER_ERROR
    ),

    /* Auth errors */
    ERR_AUTH_CANNOT_REFRESH_TOKEN("Cannot refresh token with jwt and refreshToken provided",
            ErrorOnObject.USER, ErrorOnField.USER_JWT_TOKEN, HttpStatus.UNAUTHORIZED),
    ERR_AUTH_REFRESH_TOKEN_EXPIRED("Refresh token has expired, you should sign-in again",
            ErrorOnObject.USER, ErrorOnField.USER_REFRESH_TOKEN, HttpStatus.UNAUTHORIZED),
    ERR_AUTH_JWT_TOKEN_INVALID("Cannot decode JWT token from %s",
            ErrorOnObject.USER, ErrorOnField.USER_JWT_TOKEN, HttpStatus.UNAUTHORIZED),

    /* Global errors */
    ERR_GLOBAL_SEARCH_QUERY_EMPTY("Query is required to search",
            ErrorOnObject.GLOBAL, ErrorOnField.SEARCH_QUERY, HttpStatus.BAD_REQUEST),
    ERR_GLOBAL_SEARCH_TOO_SHORT("Query provided %s is too short, should be at least 3 characters",
            ErrorOnObject.GLOBAL, ErrorOnField.SEARCH_QUERY, HttpStatus.BAD_REQUEST),
    ERR_GLOBAL_INVALID_UUID("Uuid provided %s is invalid",
            ErrorOnObject.GLOBAL, ErrorOnField.UUID, HttpStatus.BAD_REQUEST),
    ERR_GLOBAL_INVALID_PARAMETER("%s provided %s is invalid",
            ErrorOnObject.GLOBAL, null, HttpStatus.BAD_REQUEST),

    /* Internal server errors */
    ERR_UNEXPECTED_ERROR_OCCURRED("Unexpected error occurred : %s", null, null, HttpStatus.INTERNAL_SERVER_ERROR),
    ERR_CANNOT_FIND_ERROR_CODE_FROM_VALIDATION_MESSAGE("Cannot retrieve matching ErrorDetails based on validation message [%s]", null, null, HttpStatus.INTERNAL_SERVER_ERROR);

    private String message;
    private final ErrorOnObject errorOnObject;
    private final ErrorOnField errorOnField;
    private final HttpStatus httpStatus;

    ErrorDetails(String message, ErrorOnObject errorOnObject, ErrorOnField errorOnField, HttpStatus httpStatus) {
        this.message = message;
        this.errorOnObject = errorOnObject;
        this.errorOnField = errorOnField;
        this.httpStatus = httpStatus;
    }

    public ErrorDetails addingArgs(Object... args) {
        this.message = String.format(this.message, args);
        return this;
    }

    public String getCode() {
        return this.name();
    }
}
