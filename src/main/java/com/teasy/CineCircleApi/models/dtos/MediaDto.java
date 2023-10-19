package com.teasy.CineCircleApi.models.dtos;

public record MediaDto(
        Long id,
        String title,
        String originalTitle,
        String posterUrl
) {
}
