package com.teasy.CineCircleApi.models.externals.theMovieDb;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class WatchProviderInfo {
    @JsonProperty("logo_path")
    private String logoPath;
    @JsonProperty("provider_id")
    private Long providerId;
    @JsonProperty("provider_name")
    private String providerName;
    @JsonProperty("display_priority")
    private Integer displayPriority;
}
