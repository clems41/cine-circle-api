package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.annotation.Nullable;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import org.springframework.boot.context.properties.bind.DefaultValue;

public record LibraryAddMediaRequest(
        @Nullable String comment,
        @Nullable @Min(1) @Max(5) Integer rating
) {
}
