package com.teasy.CineCircleApi.models.dtos.requests;

import lombok.NonNull;

public record CircleCreateUpdateRequest(
        @NonNull String name,
        String description,
        @NonNull Boolean isPublic
) {
}
