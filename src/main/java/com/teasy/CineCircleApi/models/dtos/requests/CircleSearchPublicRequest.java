package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.annotation.Nullable;
import org.hibernate.validator.constraints.Length;
import org.springframework.web.bind.annotation.RequestParam;

public record CircleSearchPublicRequest(
        @RequestParam @Length(min = 3) String query
) {
}
