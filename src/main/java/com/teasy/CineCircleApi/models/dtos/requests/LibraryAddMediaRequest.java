package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.annotation.Nullable;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;

public record LibraryAddMediaRequest(
        @Nullable String comment,
        @Min(value = 1, message = "ERR_LIBRARY_RATING_INCORRECT")
        @Max(value = 5, message = "ERR_LIBRARY_RATING_INCORRECT") Integer rating
) {
}
