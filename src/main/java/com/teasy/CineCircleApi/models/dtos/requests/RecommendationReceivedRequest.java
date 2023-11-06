package com.teasy.CineCircleApi.models.dtos.requests;

import org.springframework.web.bind.annotation.RequestParam;

import java.util.UUID;

public record RecommendationReceivedRequest(
        @RequestParam UUID mediaId
) {
    ;
}
