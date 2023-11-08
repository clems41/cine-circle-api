package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.CircleService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/circles")
@CrossOrigin
@Tag(name = "Circle", description = "All operations related to CRUD circle + add/remove user from it")
@SecurityRequirement(name = "JWT")
public class CircleController {

    CircleService circleService;

    @Autowired
    public CircleController(CircleService circleService) {
        this.circleService = circleService;
    }

    @GetMapping("/")
    @Operation(summary = "List all circles which authenticated user belongs to")
    public ResponseEntity<?> listCircles(Principal principal) {
        try {
            return ResponseEntity.ok().body(circleService.listCircles(principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/")
    @Operation(summary = "Create new circle")
    public ResponseEntity<?> createCircle(
            @RequestBody @Valid CircleCreateUpdateRequest circleCreateUpdateRequest,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(circleService.createCircle(circleCreateUpdateRequest, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/{circle_id}")
    @Operation(summary = "Update existing circle")
    public ResponseEntity<?> updateCircle(
            @RequestBody @Valid CircleCreateUpdateRequest circleCreateUpdateRequest,
            @PathVariable("circle_id") UUID circleId,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(circleService.updateCircle(circleCreateUpdateRequest, circleId, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/{circle_id}")
    @Operation(summary = "Get existing circle")
    public ResponseEntity<?> getCircle(
            @PathVariable("circle_id") UUID circleId,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(circleService.getCircle(circleId, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @DeleteMapping("/{circle_id}")
    @Operation(summary = "Delete existing circle")
    public ResponseEntity<?> deleteCircle(
            @PathVariable("circle_id") UUID circleId,
            Principal principal
    ) {
        try {
            circleService.deleteCircle(circleId, principal.getName());
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/{circle_id}/users/{user_id}")
    @Operation(summary = "Add user to existing circle")
    public ResponseEntity<?> addUserToCircle(
            @PathVariable("circle_id") UUID circleId,
            @PathVariable("user_id") UUID userId,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(circleService.addUserToCircle(userId, circleId, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @DeleteMapping("/{circle_id}/users/{user_id}")
    @Operation(summary = "Remove user from existing circle")
    public ResponseEntity<?> removeUserFromCircle(
            @PathVariable("circle_id") UUID circleId,
            @PathVariable("user_id") UUID userId,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(circleService.removeUserFromCircle(userId, circleId, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
