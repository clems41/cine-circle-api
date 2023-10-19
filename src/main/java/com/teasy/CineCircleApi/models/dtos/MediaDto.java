package com.teasy.CineCircleApi.models.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.teasy.CineCircleApi.models.enums.MediaType;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Date;

@Getter
@Setter
@NoArgsConstructor
public class MediaDto{
    private Long id;
    private String title;
    private String originalTitle;
    private String posterUrl;
    private String backdropUrl;
    private String overview;
    private String genres;
    private String mediaType;
    @JsonFormat(pattern = "dd/MM/yyyy")
    private Date releaseDate;
    private String originalLanguage;
    private Float popularity;
    private Float voteAverage;
    private Integer voteCount;
    private String originCountry;
}
