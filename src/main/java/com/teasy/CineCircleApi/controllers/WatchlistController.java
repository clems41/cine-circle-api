package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.WatchlistService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/watchlist")
@CrossOrigin
public class WatchlistController {
    WatchlistService watchlistService;

    @Autowired
    private WatchlistController(WatchlistService watchlistService) {
        this.watchlistService = watchlistService;
    }

    @GetMapping("/")
    public ResponseEntity<?> getWatchlist(Pageable page, Principal principal) {
        try {
            return ResponseEntity.ok().body(watchlistService.getWatchlist(page, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/{mediaId}")
    public ResponseEntity<?> addToWatchlist(@PathVariable("mediaId") Long mediaId, Principal principal) {
        try {
            watchlistService.addToWatchlist(principal.getName(), mediaId);
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @DeleteMapping("/{mediaId}")
    public ResponseEntity<?> removeFromWatchlist(@PathVariable("mediaId") Long mediaId, Principal principal) {
        try {
            watchlistService.removeFromWatchlist(principal.getName(), mediaId);
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
