package com.teasy.CineCircleApi.models.dtos;

import com.teasy.CineCircleApi.models.entities.User;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Set;

@Getter
@Setter
@NoArgsConstructor
public class CircleDto {
    private String id;
    private Boolean isPublic;
    private String name;
    private String description;
    private Set<UserDto> users;
    private UserDto createdBy;
}
