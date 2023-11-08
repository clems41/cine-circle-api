package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;
import lombok.NonNull;
import org.springframework.boot.context.properties.bind.DefaultValue;
import org.springframework.web.bind.annotation.RequestParam;

public record UserSearchRequest(
        @RequestParam @NotEmpty String query
) {}
