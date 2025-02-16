package com.teasy.CineCircleApi.models.dtos.requests;

import org.springframework.web.bind.annotation.RequestParam;

public record UserSearchRelatedRequest(
        @RequestParam String query
) {}
