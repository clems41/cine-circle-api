package com.teasy.CineCircleApi.models.dtos.requests;

import org.hibernate.validator.constraints.Length;
import org.springframework.web.bind.annotation.RequestParam;

public record LibrarySearchRequest(
        @RequestParam(required = false) @Length(min = 3, message = "ERR_GLOBAL_SEARCH_TOO_SHORT") String query,
        @RequestParam(required = false) String genre
) {
}
