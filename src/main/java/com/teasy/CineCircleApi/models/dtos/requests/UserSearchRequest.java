package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;
import org.hibernate.validator.constraints.Length;
import org.springframework.web.bind.annotation.RequestParam;

public record UserSearchRequest(
        @RequestParam @NotEmpty(message = "ERR_GLOBAL_SEARCH_QUERY_EMPTY")
        @Length(min = 3, message = "ERR_GLOBAL_SEARCH_TOO_SHORT") String query
) {}
