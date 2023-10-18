package com.teasy.CineCircleApi.models.dtos.responses;

import com.teasy.CineCircleApi.models.dtos.UserDto;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class SignInResponse {
    private String token;
    private UserDto user;
}
