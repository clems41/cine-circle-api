package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;

@Getter
@Setter
public class RecommendationMediaDto {
    private String id;
    private UserDto sentBy;
    private String comment;
    private Integer rating;
    private LocalDateTime sentAt;
}
