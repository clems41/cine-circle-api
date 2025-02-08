package com.teasy.CineCircleApi.models.enums;

import lombok.Getter;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
@Getter
public enum ErrorMessage {
    /* Users errors */
    USER_NOT_FOUND("User not found"),
    USER_USERNAME_ALREADY_EXISTS("Username already exists"),
    USER_EMAIL_ALREADY_EXISTS("Email already exists"),
    USER_BAD_CREDENTIALS("Bad credentials"),
    USER_SEARCH_BAD_QUERY("Bad query for user search"),


    /* Media errors */
    MEDIA_NOT_FOUND("Media not found"),

    /* Circle errors */

    CIRCLE_NOT_FOUND("Circle not found"),
    CIRCLE_USER_BAD_PERMISSIONS("User does not have the required permissions for this circle"),

    /* Internal server errors */

    ERR_EMAILSERVICE_CANNOT_SEND_EMAIL("An error occurred while sending an email"),;

    private final String message;

    public String getCode() {
        return this.name();
    }
}
