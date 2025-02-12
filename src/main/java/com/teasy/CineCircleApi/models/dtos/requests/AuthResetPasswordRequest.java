package com.teasy.CineCircleApi.models.dtos.requests;

import org.hibernate.validator.constraints.Length;

public record AuthResetPasswordRequest(
        @Length(min = 6, message = "ERR_USER_PASSWORD_TOO_SHORT") String oldPassword,
        @Length(min = 6, message = "ERR_USER_PASSWORD_TOO_SHORT") String newPassword
) {}
