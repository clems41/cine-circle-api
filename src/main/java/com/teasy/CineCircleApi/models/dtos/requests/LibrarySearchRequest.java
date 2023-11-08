package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.annotation.Nullable;
import org.hibernate.validator.constraints.Length;
import org.springframework.boot.context.properties.bind.DefaultValue;
import org.springframework.web.bind.annotation.RequestParam;

public record LibrarySearchRequest(
        @RequestParam @Nullable @Length(min = 3) String query,
        @RequestParam @Nullable String genre
) {
}
