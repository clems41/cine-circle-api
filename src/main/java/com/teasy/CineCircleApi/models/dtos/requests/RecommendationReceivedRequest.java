package com.teasy.CineCircleApi.models.dtos.requests;

import com.teasy.CineCircleApi.utils.validators.ValidUuid;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.UUID;

public record RecommendationReceivedRequest(
        @RequestParam @ValidUuid UUID mediaId
) {
    ;
}
