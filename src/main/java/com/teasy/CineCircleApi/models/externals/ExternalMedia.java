package com.teasy.CineCircleApi.models.externals;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDate;

@Getter
@Setter
@NoArgsConstructor
public class ExternalMedia {
    private String externalId;
    private String title;
    private String originalTitle;
    private String posterUrl;
    private String backdropUrl;
    private String trailerUrl;
    private String overview;
    private String genres;
    private String mediaType;
    @JsonFormat(pattern = "dd/MM/yyyy")
    private LocalDate releaseDate;
    private Integer runtime;
    private String originalLanguage;
    private Float popularity;
    private Float voteAverage;
    private Integer voteCount;
    private String originCountry;
    private String actors;
    private String director;
}
