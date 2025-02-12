package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.Email;
import org.springframework.web.bind.annotation.RequestParam;


public record UserSendResetPasswordEmailRequest(
        @RequestParam @Email(regexp = ".+@.+\\..+", message = "ERR_USER_EMAIL_INCORRECT") String email
) {
}
