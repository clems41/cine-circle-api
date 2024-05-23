package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.Setter;

import java.util.Set;
import java.util.UUID;

@Getter
@Setter
public class UserFullInfoDto {
    private String id;
    private String username;
    private String email;
    private String displayName;
    private String topicName;
    private Set<UserDto> relatedUsers;
}
