package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;
import org.hibernate.validator.constraints.Length;
import org.springframework.web.bind.annotation.RequestParam;

public record MediaSearchRequest(
        @RequestParam @NotEmpty @Length(min = 3) String query
) {
    ;
}
