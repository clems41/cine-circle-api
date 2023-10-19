package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.SearchMediaRequest;
import com.teasy.CineCircleApi.services.HttpErrorService;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.externals.mediaProviders.TheMovieDb;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/medias")
@CrossOrigin
public class MediaController {
    MediaProvider mediaProvider;

    @Autowired
    private MediaController(TheMovieDb mediaProvider) {
        this.mediaProvider = mediaProvider;
    }

    @GetMapping("/")
    public ResponseEntity<?> searchMedias(Pageable page, SearchMediaRequest request) {
        try {
            return ResponseEntity.ok().body(mediaProvider.searchMedia(page, request));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<?> getMedia(final @PathVariable Long id) {
        try {
            return ResponseEntity.ok().body(mediaProvider.getMedia(id));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }
}
