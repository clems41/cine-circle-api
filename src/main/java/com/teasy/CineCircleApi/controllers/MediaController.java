package com.teasy.CineCircleApi.controllers;

import com.teasy.CineCircleApi.models.dtos.MediaFullDto;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.dtos.responses.MediaGenreResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.MediaService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/medias")
@CrossOrigin
@Tag(name = "Media", description = "Search and get details about medias (movies or tv shows)")
@SecurityRequirement(name = "JWT")
public class MediaController {
    MediaService mediaService;

    @Autowired
    private MediaController(MediaService mediaService) {
        this.mediaService = mediaService;
    }

    @GetMapping("")
    @Operation(summary = "Search media (movie or tv show)")
    public ResponseEntity<List<MediaShortDto>> searchMedias(
            Pageable page,
            @Valid MediaSearchRequest request
    ) {
        return ResponseEntity.ok().body(mediaService.searchMedia(page, request));
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get more details about specific media")
    public ResponseEntity<MediaFullDto> getMedia(
            @PathVariable("id") UUID id
    ) throws ExpectedException {
        return ResponseEntity.ok().body(mediaService.getMedia(id));
    }

    @GetMapping("/{id}/watch-providers")
    @Operation(summary = "Get watch providers for specific media")
    public ResponseEntity<List<String>> getWatchProviders(
            @PathVariable("id") UUID id
    ) throws ExpectedException {
        return ResponseEntity.ok().body(mediaService.getWatchProviders(id));
    }

    @GetMapping("/genres")
    @Operation(summary = "List all existing genres")
    public ResponseEntity<MediaGenreResponse> listExistingGenres() {
        return ResponseEntity.ok().body(new MediaGenreResponse(mediaService.listGenres()));
    }
}
