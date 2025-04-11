package com.teasy.CineCircleApi.models.dtos;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
public class MediaFullDto extends MediaShortDto {
    private String trailerUrl;
    private String genres;
    private Integer runtime;
    private String originalLanguage;
    private Float popularity;
    private Float voteAverage;
    private Integer voteCount;
    private String originCountry;
    private String actors;
    private String director;
}
