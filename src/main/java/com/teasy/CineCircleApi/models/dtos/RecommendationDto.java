package com.teasy.CineCircleApi.models.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;
import java.util.Set;

@Getter
@Setter
public class RecommendationDto {
    private Long id;
    private UserDto sentBy;
    private MediaDto media;
    private Set<UserDto> receivers;
    private String comment;
    private Integer rating;
    private LocalDateTime sentAt;
}
