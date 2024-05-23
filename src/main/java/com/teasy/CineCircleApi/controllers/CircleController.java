package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.CircleDto;
import com.teasy.CineCircleApi.models.dtos.CirclePublicDto;
import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.CircleSearchPublicRequest;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.CircleService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.List;
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
    public ResponseEntity<List<CircleDto>> listCircles(Principal principal) throws ExpectedException {
        return ResponseEntity.ok().body(circleService.listCircles(principal.getName()));
    }

    @GetMapping("/public")
    @Operation(summary = "Search among public circles")
    public ResponseEntity<Page<CirclePublicDto>> searchPublicCircles(
            Pageable pageable,
            @Valid CircleSearchPublicRequest circleSearchPublicRequest
    ) {
        return ResponseEntity.ok().body(circleService.searchPublicCircles(circleSearchPublicRequest, pageable));
    }

    @PutMapping("/public/{circle_id}/join")
    @Operation(summary = "Join existing public circle")
    public ResponseEntity<CirclePublicDto> joinPublicCircle(
            @PathVariable("circle_id") UUID circleId,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(circleService.joinPublicCircle(circleId, principal.getName()));
    }

    @PostMapping("/")
    @Operation(summary = "Create new circle")
    public ResponseEntity<CircleDto> createCircle(
            @RequestBody @Valid CircleCreateUpdateRequest circleCreateUpdateRequest,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(circleService.createCircle(circleCreateUpdateRequest, principal.getName()));
    }

    @PutMapping("/{circle_id}")
    @Operation(summary = "Update existing circle")
    public ResponseEntity<CircleDto> updateCircle(
            @RequestBody @Valid CircleCreateUpdateRequest circleCreateUpdateRequest,
            @PathVariable("circle_id") UUID circleId,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(circleService.updateCircle(circleCreateUpdateRequest, circleId, principal.getName()));
    }

    @GetMapping("/{circle_id}")
    @Operation(summary = "Get existing circle")
    public ResponseEntity<CircleDto> getCircle(
            @PathVariable("circle_id") UUID circleId,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(circleService.getCircle(circleId, principal.getName()));
    }

    @DeleteMapping("/{circle_id}")
    @Operation(summary = "Delete existing circle")
    public ResponseEntity<String> deleteCircle(
            @PathVariable("circle_id") UUID circleId,
            Principal principal
    ) throws ExpectedException {
        circleService.deleteCircle(circleId, principal.getName());
        return ResponseEntity.ok().body("");
    }

    @PutMapping("/{circle_id}/users/{user_id}")
    @Operation(summary = "Add user to existing circle")
    public ResponseEntity<CircleDto> addUserToCircle(
            @PathVariable("circle_id") UUID circleId,
            @PathVariable("user_id") UUID userId,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(circleService.addUserToCircle(userId, circleId, principal.getName()));
    }

    @DeleteMapping("/{circle_id}/users/{user_id}")
    @Operation(summary = "Remove user from existing circle")
    public ResponseEntity<CircleDto> removeUserFromCircle(
            @PathVariable("circle_id") UUID circleId,
            @PathVariable("user_id") UUID userId,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(circleService.removeUserFromCircle(userId, circleId, principal.getName()));
    }
}
