package com.teasy.CineCircleApi.models.dtos.responses;

import com.teasy.CineCircleApi.models.dtos.JwtTokenDto;
import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class SignInResponse {
    private JwtTokenDto token;
    private UserFullInfoDto user;
}
