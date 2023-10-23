package com.teasy.CineCircleApi.models.dtos.requests;

import lombok.NonNull;

public record UserResetPasswordRequest(
        @NonNull String email,
        @NonNull String token,
        @NonNull String newPassword
) {
}
