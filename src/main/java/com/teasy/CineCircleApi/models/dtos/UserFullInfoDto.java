package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class UserFullInfoDto {
    private Long id;
    private String username;
    private String email;
    private String displayName;
}
