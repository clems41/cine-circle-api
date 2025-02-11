package com.teasy.CineCircleApi.models.dtos.requests;

import com.teasy.CineCircleApi.models.enums.RecommendationType;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.UUID;

public record RecommendationSearchRequest(
        @RequestParam(required = false) UUID mediaId,
        @RequestParam(required = false) RecommendationType type,
        @RequestParam(required = false) Boolean read
) {
}
