package com.teasy.CineCircleApi.models.dtos.requests;

import com.teasy.CineCircleApi.utils.validators.ValidUuid;
import jakarta.annotation.Nullable;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.Size;

import java.util.List;
import java.util.UUID;

public record RecommendationCreateRequest(
    @ValidUuid UUID mediaId,
    @Size(min = 1) List<UUID> userIds,
    @Nullable String comment,
    @Nullable @Min(1) @Max(5) Integer rating
){}
