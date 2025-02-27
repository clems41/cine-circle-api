package com.teasy.CineCircleApi.services.utils;

import org.springframework.data.domain.ExampleMatcher;

public abstract class CustomExampleMatcher {
    public static ExampleMatcher matchingAll() {
        return ExampleMatcher.matchingAll()
                .withIgnorePaths("createdAt", "updatedAt");
    }
    public static ExampleMatcher matchingAny() {
        return ExampleMatcher.matchingAny()
                .withIgnorePaths("createdAt", "updatedAt");
    }
}
