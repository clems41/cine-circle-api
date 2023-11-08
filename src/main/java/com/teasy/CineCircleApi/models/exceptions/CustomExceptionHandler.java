package com.teasy.CineCircleApi.models.exceptions;

import com.teasy.CineCircleApi.models.enums.ErrorCodeEnum;
import org.springframework.http.HttpStatus;

import java.util.UUID;

public abstract class CustomExceptionHandler {
    public static CustomException userWithUsernameNotFound(String username) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCodeEnum.USER_NOT_FOUND,
                String.format("user with username %s cannot be found", username)
        );
    }

    public static CustomException userWithIdNotFound(UUID id) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCodeEnum.USER_NOT_FOUND,
                String.format("user with id %s cannot be found", id)
        );
    }

    public static CustomException userWithEmailNotFound(String email) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCodeEnum.USER_NOT_FOUND,
                String.format("user with email %s cannot be found", email)
        );
    }

    public static CustomException userWithUsernameOrEmailNotFound(String username, String email) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCodeEnum.USER_NOT_FOUND,
                String.format("user with username %s or email %s cannot be found", username, email)
        );
    }
    public static CustomException userWithUsernameAlreadyExists(String username) {
        return new CustomException(
                HttpStatus.BAD_REQUEST,
                ErrorCodeEnum.USER_ALREADY_EXISTS,
                String.format("user with username %s already exists", username)
        );
    }
    public static CustomException userWithEmailAlreadyExists(String email) {
        return new CustomException(
                HttpStatus.BAD_REQUEST,
                ErrorCodeEnum.USER_ALREADY_EXISTS,
                String.format("user with email %s already exists", email)
        );
    }
    public static CustomException userWithUsernameBadCredentials(String username) {
        return new CustomException(
                HttpStatus.FORBIDDEN,
                ErrorCodeEnum.USER_BAD_CREDENTIALS,
                String.format("wrong password or token given for user with username %s", username)
        );
    }

    public static CustomException userWithEmailBadCredentials(String email) {
        return new CustomException(
                HttpStatus.FORBIDDEN,
                ErrorCodeEnum.USER_BAD_CREDENTIALS,
                String.format("wrong password or token given for user with email %s", email)
        );
    }
    public static CustomException userSearchQueryEmpty() {
        return new CustomException(
                HttpStatus.BAD_REQUEST,
                ErrorCodeEnum.USER_SEARCH_BAD_QUERY,
                "query is empty"
        );
    }






    public static CustomException mediaWithIdNotFound(UUID id) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCodeEnum.MEDIA_NOT_FOUND,
                String.format("media with id %s cannot be found", id)
        );
    }
    public static CustomException mediaWithExternalIdNotFound(String externalId) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCodeEnum.MEDIA_NOT_FOUND,
                String.format("media with externalId %s cannot be found", externalId)
        );
    }




    public static CustomException circleWithIdNotFound(UUID id) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCodeEnum.CIRCLE_NOT_FOUND,
                String.format("circle with id %s cannot be found", id)
        );
    }

    public static CustomException circleWithIdUserWithUsernameBadPermissions(UUID circleId, String username) {
        return new CustomException(
                HttpStatus.FORBIDDEN,
                ErrorCodeEnum.CIRCLE_USER_BAD_PERMISSIONS,
                String.format("circle with id %s cannot be updated/deleted by user with username %s",
                        circleId, username)
        );
    }

}
