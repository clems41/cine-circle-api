package com.teasy.CineCircleApi.utils;

import lombok.Getter;
import org.springframework.web.util.UriComponentsBuilder;

import java.util.Map;

@Getter
public abstract class HttpUtils {
    public final static String authUrl = "/auth/";
    public final static String userUrl = "/users/";
    public final static String libraryUrl = "/library/";
    public final static String circleUrl = "/circles/";
    public final static String mediaUrl = "/medias/";
    public final static String recommendationUrl = "/recommendations";
    public final static String watchlistUrl = "/watchlist";
    public static String getTestingUrl(int port) {
        return "http://localhost:".concat(String.valueOf(port)).concat("/api/v1");
    }

    public static String getUriWithQueryParameter(int port, String path, Map<String, Object> queryParameters) {
        var uriBuilder = UriComponentsBuilder.newInstance()
                .scheme("http")
                .host("localhost:".concat(String.valueOf(port)))
                .path("/api/v1".concat(path));
        for (Map.Entry<String, Object> entry : queryParameters.entrySet()) {
            uriBuilder.queryParam(entry.getKey(), entry.getValue());
        }
        return uriBuilder.toUriString().replace("%3A", ":");
    }
}
