package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.dtos.responses.MediaGenreResponse;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.MediaService;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
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

    @GetMapping("/")
    @Operation(summary = "Search media (movie or tv show)")
    public ResponseEntity<?> searchMedias(Pageable page, MediaSearchRequest request, Principal principal) {
        try {
            return ResponseEntity.ok().body(mediaService.searchMedia(page, request, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get more details about specific media")
    public ResponseEntity<?> getMedia(final @PathVariable("id") UUID id, Principal principal) {
        try {
            return ResponseEntity.ok().body(mediaService.getMedia(id, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/genres")
    @Operation(summary = "List all existing genres")
    public ResponseEntity<?> listExistingGenres() {
        try {
            return ResponseEntity.ok().body(new MediaGenreResponse(mediaService.listGenres()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
