package com.teasy.CineCircleApi.models.externals.theMovieDb;

import lombok.Getter;
import lombok.Setter;

import java.util.List;

@Getter
@Setter
public class GenreMovieListResponse {
    private List<Genre> genres;
}
