package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.annotation.Nullable;
import jakarta.validation.constraints.NotEmpty;
import lombok.NonNull;
import org.springframework.boot.context.properties.bind.DefaultValue;

public record CircleCreateUpdateRequest(
        @NotEmpty String name,
        @DefaultValue("") String description,
        Boolean isPublic
) {
}
