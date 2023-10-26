package com.teasy.CineCircleApi.models.dtos.requests;

import org.springframework.web.bind.annotation.RequestParam;

public record RecommendationReceivedRequest(
        @RequestParam Long mediaId
) {
    ;
}
