package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;
import org.springframework.boot.context.properties.bind.DefaultValue;

public record CircleCreateUpdateRequest(
        @NotEmpty(message = "ERR_CIRCLE_NAME_EMPTY") String name,
        @DefaultValue("") String description,
        @DefaultValue("false") Boolean isPublic
) {
}
