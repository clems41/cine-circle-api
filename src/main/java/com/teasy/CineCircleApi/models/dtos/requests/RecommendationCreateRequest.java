package com.teasy.CineCircleApi.models.dtos.requests;

import lombok.NonNull;
import java.util.List;

public record RecommendationCreateRequest(
    @NonNull Long mediaId,
    @NonNull List<Long> userIds,
    @NonNull String comment,
    @NonNull Integer rating
){}
