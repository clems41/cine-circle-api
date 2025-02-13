package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;

public record AuthRefreshTokenRequest(
        @NotEmpty(message = "ERR_USER_JWT_TOKEN_EMPTY") String jwtToken,
        @NotEmpty(message = "ERR_USER_REFRESH_TOKEN_EMPTY") String jwtRefreshToken
) {}
