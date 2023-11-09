package com.teasy.CineCircleApi.models.externals.theMovieDb;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class WatchProvidersResponse {
    private Long id;
    private WatchProvidersResults results;
}
