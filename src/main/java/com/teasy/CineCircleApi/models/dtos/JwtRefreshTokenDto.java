package com.teasy.CineCircleApi.models.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;

import java.time.LocalDateTime;

public record JwtRefreshTokenDto(
        String tokenString,
        @JsonFormat(timezone = "Europe/Paris")
        LocalDateTime expirationDate
) {
}
