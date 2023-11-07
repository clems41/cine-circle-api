package com.teasy.CineCircleApi.utils;

public abstract class HttpUtils {
    public static String getTestingUrl(int port) {
        return "http://localhost:".concat(String.valueOf(port)).concat("/api/v1");
    }
}
