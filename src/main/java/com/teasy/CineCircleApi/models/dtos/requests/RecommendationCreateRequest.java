package com.teasy.CineCircleApi.models.dtos.requests;

import lombok.NonNull;
import java.util.List;
import java.util.UUID;

public record RecommendationCreateRequest(
    @NonNull UUID mediaId,
    @NonNull List<UUID> userIds,
    @NonNull String comment,
    @NonNull Integer rating
){}
