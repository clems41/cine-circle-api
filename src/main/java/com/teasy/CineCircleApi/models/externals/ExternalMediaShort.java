package com.teasy.CineCircleApi.models.externals;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDate;

@Getter
@Setter
public class ExternalMediaShort {
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
}
