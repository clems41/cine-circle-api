package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.services.CircleService;
import com.teasy.CineCircleApi.services.HttpErrorService;
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
@RequestMapping("/api/v1/circles")
@CrossOrigin
public class CircleController {

    CircleService circleService;

    @Autowired
    public CircleController(CircleService circleService) {
        this.circleService = circleService;
    }

    @GetMapping("/")
    public ResponseEntity<?> listCircles() {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(circleService.listCircles(usernameFromToken));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @PostMapping("/")
    public ResponseEntity<?> createCircle(@RequestBody CircleCreateUpdateRequest circleCreateUpdateRequest) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(circleService.createCircle(circleCreateUpdateRequest, usernameFromToken));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @PutMapping("/{circle_id}")
    public ResponseEntity<?> createCircle(@RequestBody CircleCreateUpdateRequest circleCreateUpdateRequest,
                                          @PathVariable("circle_id") Long circleId) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(circleService.updateCircle(circleCreateUpdateRequest, circleId, usernameFromToken));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @DeleteMapping("/{circle_id}")
    public ResponseEntity<?> createCircle(@PathVariable("circle_id") Long circleId) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            circleService.deleteCircle(circleId, usernameFromToken);
            return ResponseEntity.ok().body("");
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @PutMapping("/{circle_id}/users/{user_id}")
    public ResponseEntity<?> addUserToCircle(@PathVariable("circle_id") Long circleId,
                                             @PathVariable("user_id") Long userId) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(circleService.addUserToCircle(userId, circleId, usernameFromToken));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @DeleteMapping("/{circle_id}/users/{user_id}")
    public ResponseEntity<?> removeUserFromCircle(@PathVariable("circle_id") Long circleId,
                                                  @PathVariable("user_id") Long userId) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(circleService.removeUserFromCircle(userId, circleId, usernameFromToken));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }
}
