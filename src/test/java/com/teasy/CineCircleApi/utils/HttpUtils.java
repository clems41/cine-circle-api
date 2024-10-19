package com.teasy.CineCircleApi.utils;

import lombok.Getter;

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
}
