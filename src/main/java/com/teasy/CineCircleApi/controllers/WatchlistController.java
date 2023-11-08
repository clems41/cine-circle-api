package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.WatchlistService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/watchlist")
@CrossOrigin
@Tag(name = "Watchlist", description = "Add, remove and list medias from authenticated user watchlist")
@SecurityRequirement(name = "JWT")
public class WatchlistController {
    WatchlistService watchlistService;

    @Autowired
    private WatchlistController(WatchlistService watchlistService) {
        this.watchlistService = watchlistService;
    }

    @GetMapping("/")
    @Operation(summary = "List medias from authenticated user watchlist")
    public ResponseEntity<?> getWatchlist(
            Pageable page,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(watchlistService.getWatchlist(page, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/{mediaId}")
    @Operation(summary = "Add media to authenticated user watchlist")
    public ResponseEntity<?> addToWatchlist(
            @PathVariable("mediaId") UUID mediaId,
            Principal principal
    ) {
        try {
            watchlistService.addToWatchlist(principal.getName(), mediaId);
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @DeleteMapping("/{mediaId}")
    @Operation(summary = "Remove media from authenticated user watchlist")
    public ResponseEntity<?> removeFromWatchlist(
            @PathVariable("mediaId") UUID mediaId,
            Principal principal
    ) {
        try {
            watchlistService.removeFromWatchlist(principal.getName(), mediaId);
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
