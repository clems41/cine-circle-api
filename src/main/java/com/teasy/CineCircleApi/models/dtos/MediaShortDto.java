package com.teasy.CineCircleApi.models.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDate;

@Getter
@Setter
@NoArgsConstructor
public class MediaShortDto {
    private String id;
    private String title;
    private String originalTitle;
    private String posterUrl;
    private String backdropUrl;
    private String overview;
    private MediaTypeEnum mediaType;
    @JsonFormat(pattern = "dd/MM/yyyy")
    private LocalDate releaseDate;
    private Integer recommendationRatingCount;
    private Double recommendationRatingAverage;
    private String personalComment;
    private Integer personalRating;
}
