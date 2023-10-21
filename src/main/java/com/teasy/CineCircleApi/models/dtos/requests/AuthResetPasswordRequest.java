package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.Email;
import lombok.NonNull;
import org.hibernate.validator.constraints.Length;

public record AuthResetPasswordRequest(
        @NonNull @Length(min = 6) String oldPassword,
        @NonNull @Length(min = 6) String newPassword
) {}
