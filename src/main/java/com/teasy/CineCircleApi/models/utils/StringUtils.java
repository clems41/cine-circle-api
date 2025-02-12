package com.teasy.CineCircleApi.models.utils;

public abstract class StringUtils {
    private final static int MAX_CHAR_POSTGRESQL = 250;
    public static String substringForDatabase(String str) {
        return str.substring(0, Math.min(str.length(), MAX_CHAR_POSTGRESQL)); // max char in database
    }
}
