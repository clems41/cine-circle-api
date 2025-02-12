package com.teasy.CineCircleApi.models.dtos.requests;

import org.hibernate.validator.constraints.Length;
import org.springframework.web.bind.annotation.RequestParam;

public record CircleSearchPublicRequest(
        @RequestParam @Length(min = 3, message = "ERR_GLOBAL_SEARCH_QUERY_EMPTY") String query
) {
}
