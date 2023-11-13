package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Set;

@Getter
@Setter
@NoArgsConstructor
public class CirclePublicDto {
    private String id;
    private Boolean isPublic;
    private String name;
    private String description;
    private UserDto createdBy;
}
