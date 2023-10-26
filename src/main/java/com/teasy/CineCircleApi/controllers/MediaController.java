package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.externals.mediaProviders.theMovieDb.TheMovieDbService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/medias")
@CrossOrigin
public class MediaController {
    MediaProvider mediaProvider;

    @Autowired
    private MediaController(TheMovieDbService mediaProvider) {
        this.mediaProvider = mediaProvider;
    }

    @GetMapping("/")
    public ResponseEntity<?> searchMedias(Pageable page, MediaSearchRequest request) {
        try {
            return ResponseEntity.ok().body(mediaProvider.searchMedia(page, request));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<?> getMedia(final @PathVariable Long id) {
        try {
            return ResponseEntity.ok().body(mediaProvider.getMedia(id));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
