package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.LibraryAddMediaRequest;
import com.teasy.CineCircleApi.models.dtos.requests.LibrarySearchRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.LibraryService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/library")
@CrossOrigin
@Tag(name = "Library", description = "Add, remove and list medias from authenticated user library")
@SecurityRequirement(name = "JWT")
public class LibraryController {
    LibraryService libraryService;

    @Autowired
    private LibraryController(LibraryService libraryService) {
        this.libraryService = libraryService;
    }

    @GetMapping("/")
    @Operation(summary = "Search medias among authenticated user library")
    public ResponseEntity<?> searchInLibrary(
            Pageable page,
            @Valid LibrarySearchRequest librarySearchRequest,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(libraryService.searchInLibrary(page, librarySearchRequest, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/{mediaId}")
    @Operation(summary = "Add media to authenticated user library")
    public ResponseEntity<?> addToLibrary(
            @PathVariable("mediaId") UUID mediaId,
            @Valid @RequestBody LibraryAddMediaRequest libraryAddMediaRequest,
            Principal principal
    ) {
        try {
            libraryService.addToLibrary(mediaId, libraryAddMediaRequest, principal.getName());
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @DeleteMapping("/{mediaId}")
    @Operation(summary = "Remove media from authenticated user library")
    public ResponseEntity<?> removeFromLibrary(
            @PathVariable("mediaId") UUID mediaId,
            Principal principal
    ) {
        try {
            libraryService.removeFromLibrary(principal.getName(), mediaId);
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
