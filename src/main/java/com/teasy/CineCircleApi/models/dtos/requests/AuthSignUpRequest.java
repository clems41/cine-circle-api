package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.Email;
import org.hibernate.validator.constraints.Length;

public record AuthSignUpRequest(
        @Length(min = 4, message = "ERR_USER_USERNAME_TOO_SHORT") String username,
        @Email(regexp = ".+@.+\\..+", message = "ERR_USER_EMAIL_INCORRECT") String email,
        @Length(min = 6, message = "ERR_USER_PASSWORD_TOO_SHORT") String password,
        @Length(min = 4, max = 20, message = "ERR_USER_DISPLAY_NAME_INCORRECT") String displayName
) {}
