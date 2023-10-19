package com.teasy.CineCircleApi.models.dtos.requests;

import org.springframework.web.bind.annotation.RequestParam;

public record SearchMediaRequest(
        @RequestParam String query
) {
    ;
}
