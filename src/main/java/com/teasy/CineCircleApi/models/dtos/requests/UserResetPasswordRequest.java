package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotEmpty;
import org.hibernate.validator.constraints.Length;

public record UserResetPasswordRequest(
        @Email(regexp = ".+@.+\\..+", message = "ERR_USER_EMAIL_INCORRECT") String email,
        @NotEmpty(message = "ERR_USER_RESET_PASSWORD_TOKEN_EMPTY") String token,
        @Length(min = 6, message = "ERR_USER_PASSWORD_TOO_SHORT") String newPassword
) {
}
