package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;
import java.util.Set;

@Getter
@Setter
public class MediaRecommendationDto {
    private Long id;
    private UserDto sentBy;
    private String comment;
    private Integer rating;
    private LocalDateTime sentAt;
}
