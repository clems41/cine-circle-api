package com.teasy.CineCircleApi.models.externals.theMovieDb;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;

import java.util.Date;
import java.util.List;

@Getter
@Setter
public class SearchResult {
    private Long id;
    @JsonProperty("poster_path")
    private String posterPath;
    @JsonProperty("media_type")
    private String mediaType;
    private String title;
    @JsonProperty("original_title")
    private String originalTitle;
    private String name;
    @JsonProperty("original_name")
    private String originalName;
    private String overview;
    private Boolean adult;
    @JsonProperty("backdrop_path")
    private String backdropPath;
    @JsonProperty("original_language")
    private String originalLanguage;
    @JsonProperty("genre_ids")
    private List<Long> genreIds;
    private Float popularity;
    @JsonProperty("first_air_date")
    private Date firstAirDate;
    @JsonProperty("release_date")
    private Date releaseDate;
    @JsonProperty("vote_average")
    private Float voteAverage;
    @JsonProperty("vote_count")
    private Integer voteCount;
    private Boolean video;
    @JsonProperty("origin_country")
    private List<String> originCountry;
}
