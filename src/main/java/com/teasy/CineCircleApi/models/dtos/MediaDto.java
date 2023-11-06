package com.teasy.CineCircleApi.models.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDate;

@Getter
@Setter
@NoArgsConstructor
public class MediaDto{
    private String id;
    private String title;
    private String originalTitle;
    private String posterUrl;
    private String backdropUrl;
    private String overview;
    private String mediaType;
    @JsonFormat(pattern = "dd/MM/yyyy")
    private LocalDate releaseDate;
    private Integer runtime;
    private String originalLanguage;
    private Integer recommendationRatingCount;
    private Double recommendationRatingAverage;
}
