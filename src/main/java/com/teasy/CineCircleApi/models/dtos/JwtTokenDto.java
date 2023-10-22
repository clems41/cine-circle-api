package com.teasy.CineCircleApi.models.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDate;
import java.util.Date;

public record JwtTokenDto(
        String tokenString,
        @JsonFormat(timezone = "Pacific/Noumea")
        Date expirationDate
) {
}
