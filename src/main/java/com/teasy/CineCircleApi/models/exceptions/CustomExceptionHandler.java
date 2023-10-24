package com.teasy.CineCircleApi.models.exceptions;

import com.teasy.CineCircleApi.models.enums.ErrorCode;
import org.springframework.http.HttpStatus;

public abstract class CustomExceptionHandler {
    public static CustomException userWithUsernameNotFound(String username) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCode.USER_NOT_FOUND,
                String.format("user with username %s cannot be found", username)
        );
    }

    public static CustomException userWithIdNotFound(Long id) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCode.USER_NOT_FOUND,
                String.format("user with id %s cannot be found", id)
        );
    }

    public static CustomException userWithEmailNotFound(String email) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCode.USER_NOT_FOUND,
                String.format("user with email %s cannot be found", email)
        );
    }

    public static CustomException userWithUsernameOrEmailNotFound(String username, String email) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCode.USER_NOT_FOUND,
                String.format("user with username %s or email %s cannot be found", username, email)
        );
    }
    public static CustomException userWithUsernameAlreadyExists(String username) {
        return new CustomException(
                HttpStatus.BAD_REQUEST,
                ErrorCode.USER_ALREADY_EXISTS,
                String.format("user with username %s already exists", username)
        );
    }
    public static CustomException userWithEmailAlreadyExists(String email) {
        return new CustomException(
                HttpStatus.BAD_REQUEST,
                ErrorCode.USER_ALREADY_EXISTS,
                String.format("user with email %s already exists", email)
        );
    }
    public static CustomException userWithUsernameBadCredentials(String username) {
        return new CustomException(
                HttpStatus.FORBIDDEN,
                ErrorCode.USER_BAD_CREDENTIALS,
                String.format("wrong password or token given for user with username %s", username)
        );
    }

    public static CustomException userWithEmailBadCredentials(String email) {
        return new CustomException(
                HttpStatus.FORBIDDEN,
                ErrorCode.USER_BAD_CREDENTIALS,
                String.format("wrong password or token given for user with email %s", email)
        );
    }
    public static CustomException userSearchQueryEmpty() {
        return new CustomException(
                HttpStatus.BAD_REQUEST,
                ErrorCode.USER_SEARCH_BAD_QUERY,
                "query is empty"
        );
    }






    public static CustomException mediaWithIdNotFound(Long id) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCode.MEDIA_NOT_FOUND,
                String.format("media with id %d cannot be found", id)
        );
    }




    public static CustomException circleWithIdNotFound(Long id) {
        return new CustomException(
                HttpStatus.NOT_FOUND,
                ErrorCode.CIRCLE_NOT_FOUND,
                String.format("circle with id %d cannot be found", id)
        );
    }

    public static CustomException circleWithIdUserWithUsernameBadPermissions(Long circleId, String username) {
        return new CustomException(
                HttpStatus.FORBIDDEN,
                ErrorCode.CIRCLE_USER_BAD_PERMISSIONS,
                String.format("circle with id %d cannot be updated/deleted by user with username %s",
                        circleId, username)
        );
    }

}
