package com.teasy.CineCircleApi.models.externals.theMovieDb;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class WatchProvidersResults {
    @JsonProperty("FR")
    private WatchProvidersByRegion fr;
}
