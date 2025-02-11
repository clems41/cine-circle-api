package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationSearchRequest;
import com.teasy.CineCircleApi.models.dtos.responses.RecommendationCreateResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.RecommendationService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/recommendations")
@CrossOrigin
@Tag(name = "Recommendation", description = "Send and list recommendations")
@SecurityRequirement(name = "JWT")
public class RecommendationController {

    RecommendationService recommendationService;

    @Autowired
    public RecommendationController(RecommendationService recommendationService) {
        this.recommendationService = recommendationService;
    }

    @GetMapping("")
    @Operation(summary = "Search among all recommendations (received and sent)")
    public ResponseEntity<Page<RecommendationDto>> listReceivedRecommendations(
            @PageableDefault(sort = "sentAt", direction = Sort.Direction.DESC) Pageable pageable,
            @Valid RecommendationSearchRequest recommendationSearchRequest,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(recommendationService.searchRecommendations(
                pageable,
                recommendationSearchRequest,
                principal.getName())
        );
    }

    @PutMapping("/{recommendationId}/read")
    @Operation(summary = "Mark a specific recommendation as read for the authenticated user")
    public ResponseEntity<String> markRecommendationAsRead(
            @PathVariable("recommendationId") UUID recommendationId,
            Principal principal
    ) throws ExpectedException {
        recommendationService.markRecommendationAsRead(recommendationId, principal.getName());
        return ResponseEntity.ok().body("");
    }

    @PostMapping("")
    @Operation(summary = "Send new recommendation")
    public ResponseEntity<RecommendationCreateResponse> createRecommendation(
            @RequestBody @Valid RecommendationCreateRequest recommendationCreateRequest,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(recommendationService.createRecommendation(
                recommendationCreateRequest,
                principal.getName())
        );
    }
}
