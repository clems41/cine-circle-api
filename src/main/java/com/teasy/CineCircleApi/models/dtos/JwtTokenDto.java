package com.teasy.CineCircleApi.models.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;

import java.util.Date;

public record JwtTokenDto(
        String tokenString,
        @JsonFormat(timezone = "Europe/Paris")
        Date expirationDate
) {
}
