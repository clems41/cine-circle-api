package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.LibraryService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

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
        } catch (CustomException e) {
            return e.getEntityResponse();
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
        } catch (CustomException e) {
            return e.getEntityResponse();
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
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
