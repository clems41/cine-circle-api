package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.CircleService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;

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
    public ResponseEntity<?> listCircles(Principal principal) {
        try {
            return ResponseEntity.ok().body(circleService.listCircles(principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/")
    public ResponseEntity<?> createCircle(@RequestBody CircleCreateUpdateRequest circleCreateUpdateRequest,
                                          Principal principal) {
        try {
            return ResponseEntity.ok().body(circleService.createCircle(circleCreateUpdateRequest, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/{circle_id}")
    public ResponseEntity<?> createCircle(@RequestBody CircleCreateUpdateRequest circleCreateUpdateRequest,
                                          @PathVariable("circle_id") Long circleId,
                                          Principal principal) {
        try {
            return ResponseEntity.ok().body(circleService.updateCircle(circleCreateUpdateRequest, circleId, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @DeleteMapping("/{circle_id}")
    public ResponseEntity<?> createCircle(@PathVariable("circle_id") Long circleId, Principal principal) {
        try {
            circleService.deleteCircle(circleId, principal.getName());
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/{circle_id}/users/{user_id}")
    public ResponseEntity<?> addUserToCircle(@PathVariable("circle_id") Long circleId,
                                             @PathVariable("user_id") Long userId,
                                             Principal principal) {
        try {
            return ResponseEntity.ok().body(circleService.addUserToCircle(userId, circleId, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @DeleteMapping("/{circle_id}/users/{user_id}")
    public ResponseEntity<?> removeUserFromCircle(@PathVariable("circle_id") Long circleId,
                                                  @PathVariable("user_id") Long userId,
                                                  Principal principal) {
        try {
            return ResponseEntity.ok().body(circleService.removeUserFromCircle(userId, circleId, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
