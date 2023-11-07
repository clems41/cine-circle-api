package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;
import java.util.Set;

@Getter
@Setter
public class RecommendationDto {
    private String id;
    private UserDto sentBy;
    private MediaShortDto media;
    private Set<UserDto> receivers;
    private String comment;
    private Integer rating;
    private LocalDateTime sentAt;
}
