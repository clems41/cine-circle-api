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
    ERR_USER_PASSWORD_INCORRECT("Password provided for user with id %s does not match the one in database",
            ErrorOnObject.USER, ErrorOnField.USER_PASSWORD, HttpStatus.BAD_REQUEST),
    ERR_USER_RESET_PASSWORD_TOKEN_INCORRECT("Reset password token provided for user with id %s does not match the one in database",
            ErrorOnObject.USER, ErrorOnField.USER_RESET_PASSWORD_TOKEN, HttpStatus.BAD_REQUEST),
    ERR_USER_SEARCH_QUERY_EMPTY("Query to search users cannot be empty",
            ErrorOnObject.USER, ErrorOnField.USER_SEARCH_QUERY, HttpStatus.BAD_REQUEST),

    /* Library errors */
    ERR_LIBRARY_NOT_FOUND("Library for username %s and mediaId %S cannot be found",
            ErrorOnObject.LIBRARY, null, HttpStatus.NOT_FOUND),

    /* Heading errors */

    /* Contact errors */

    /* Notification errors */

    /* Recommendation errors */

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

    /* Email errors */
    ERR_EMAIL_BUILDING_REQUEST(
            "An error occurred while building the email from request",
            ErrorOnObject.EMAIL_SERVICE,
            null,
            HttpStatus.INTERNAL_SERVER_ERROR
    ),
    ERR_EMAIL_SENDING_REQUEST(
            "An error occurred while sending the email to %s",
            ErrorOnObject.EMAIL_SERVICE,
            null,
            HttpStatus.INTERNAL_SERVER_ERROR
    ),

    /* Internal server errors */
    ERR_UNEXPECTED_ERROR_OCCURRED("Unexpected error occurred : %s", null, null, HttpStatus.INTERNAL_SERVER_ERROR);

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

    public boolean isSameErrorThan(ErrorDetails object) {
        return this.getCode().equals(object.getCode());
    }
}
