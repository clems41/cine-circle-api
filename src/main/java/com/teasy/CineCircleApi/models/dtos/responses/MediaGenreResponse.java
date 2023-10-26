package com.teasy.CineCircleApi.models.dtos.responses;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.List;

@Getter
@AllArgsConstructor
public class MediaGenreResponse {
    private List<String> genres;
}
