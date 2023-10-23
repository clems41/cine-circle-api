package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.services.HttpErrorService;
import com.teasy.CineCircleApi.services.LibraryService;
import com.teasy.CineCircleApi.services.WatchlistService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/library")
@CrossOrigin
public class LibraryController {
    LibraryService libraryService;

    @Autowired
    private LibraryController(LibraryService libraryService) {
        this.libraryService = libraryService;
    }

    @GetMapping("/")
    public ResponseEntity<?> getLibrary(Pageable page) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(libraryService.getLibrary(page, usernameFromToken));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @PutMapping("/{mediaId}")
    public ResponseEntity<?> addToLibrary(@PathVariable("mediaId") Long mediaId) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            libraryService.addToLibrary(usernameFromToken, mediaId);
            return ResponseEntity.ok().body("");
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @DeleteMapping("/{mediaId}")
    public ResponseEntity<?> removeFromLibrary(@PathVariable("mediaId") Long mediaId) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            libraryService.removeFromLibrary(usernameFromToken, mediaId);
            return ResponseEntity.ok().body("");
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }
}
