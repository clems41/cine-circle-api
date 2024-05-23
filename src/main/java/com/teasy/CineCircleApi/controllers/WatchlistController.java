package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.WatchlistService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
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
    public ResponseEntity<Page<MediaShortDto>> getWatchlist(
            Pageable page,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(watchlistService.getWatchlist(page, principal.getName()));
    }

    @PutMapping("/{mediaId}")
    @Operation(summary = "Add media to authenticated user watchlist")
    public ResponseEntity<String> addToWatchlist(
            @PathVariable("mediaId") UUID mediaId,
            Principal principal
    ) throws ExpectedException {
        watchlistService.addToWatchlist(principal.getName(), mediaId);
        return ResponseEntity.ok().body("");
    }

    @DeleteMapping("/{mediaId}")
    @Operation(summary = "Remove media from authenticated user watchlist")
    public ResponseEntity<String> removeFromWatchlist(
            @PathVariable("mediaId") UUID mediaId,
            Principal principal
    ) throws ExpectedException {
        watchlistService.removeFromWatchlist(principal.getName(), mediaId);
        return ResponseEntity.ok().body("");
    }
}
