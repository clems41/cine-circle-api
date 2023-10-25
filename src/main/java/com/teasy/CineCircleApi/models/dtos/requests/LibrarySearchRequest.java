package com.teasy.CineCircleApi.models.dtos.requests;

import org.springframework.boot.context.properties.bind.DefaultValue;
import org.springframework.web.bind.annotation.RequestParam;

public record LibrarySearchRequest(
        @RequestParam @DefaultValue(value = "") String query
) {
}
