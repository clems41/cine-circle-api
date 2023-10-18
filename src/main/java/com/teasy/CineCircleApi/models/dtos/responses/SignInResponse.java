package com.teasy.CineCircleApi.models.dtos.responses;

import com.teasy.CineCircleApi.models.dtos.UserDetails;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class SignInResponse {
    private String token;
    private UserDetails user;
}
