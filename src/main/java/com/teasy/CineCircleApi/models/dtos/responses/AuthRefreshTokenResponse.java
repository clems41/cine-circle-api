package com.teasy.CineCircleApi.models.dtos.responses;

import com.teasy.CineCircleApi.models.dtos.JwtTokenDto;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class AuthRefreshTokenResponse {
    private JwtTokenDto jwtToken;
}
