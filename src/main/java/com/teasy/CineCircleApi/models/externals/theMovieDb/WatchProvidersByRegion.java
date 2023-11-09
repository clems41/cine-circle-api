package com.teasy.CineCircleApi.models.externals.theMovieDb;

import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import lombok.Getter;
import lombok.Setter;

import java.util.List;

@Getter
@Setter
public class WatchProvidersByRegion {
    private String link;
    private List<WatchProviderInfo> buy;
    private List<WatchProviderInfo> rent;
    private List<WatchProviderInfo> flatrate;
}
