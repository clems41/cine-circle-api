package com.teasy.CineCircleApi.models.dtos.requests;

import org.springframework.web.bind.annotation.RequestParam;

import java.util.UUID;

public record RecommendationSearchRequest(
        @RequestParam(required = false) UUID mediaId,

        @RequestParam(required = false) String type,

        @RequestParam(required = false) Boolean read
) {
}
