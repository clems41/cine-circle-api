package com.teasy.CineCircleApi.models.dtos.requests;

import com.teasy.CineCircleApi.utils.validators.ValidUuid;
import jakarta.annotation.Nullable;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.Size;

import java.util.List;
import java.util.UUID;

public record RecommendationCreateRequest(
    @ValidUuid(message = "ERR_RECOMMENDATION_MEDIA_ID_INCORRECT") UUID mediaId,
    @Size(min = 1, message = "ERR_RECOMMENDATION_USER_IDS_INCORRECT") List<UUID> userIds,
    @Nullable List<UUID> circleIds,
    @Nullable String comment,
    @Nullable @Min(value = 1, message = "ERR_RECOMMENDATION_RATING_INCORRECT")
    @Max(value =5, message = "ERR_RECOMMENDATION_RATING_INCORRECT") Integer rating
){}
