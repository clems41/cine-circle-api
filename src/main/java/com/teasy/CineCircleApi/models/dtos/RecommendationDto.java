package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;

@Getter
@Setter
public class RecommendationDto {
    private String id;
    private String recommendationRef;
    private UserDto sentBy;
    private MediaShortDto media;
    private UserDto receiver;
    private String comment;
    private Integer rating;
    private LocalDateTime sentAt;
}
