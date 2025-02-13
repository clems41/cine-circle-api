package com.teasy.CineCircleApi.models.dtos.responses;

import com.teasy.CineCircleApi.models.dtos.JwtRefreshTokenDto;
import com.teasy.CineCircleApi.models.dtos.JwtTokenDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class AuthSignInResponse {
    private JwtTokenDto token;
    private JwtRefreshTokenDto refreshToken;
    private UserFullInfoDto user;
}
