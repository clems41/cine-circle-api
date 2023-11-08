package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationReceivedRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.RecommendationService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;

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

    @GetMapping("/received")
    @Operation(summary = "Search among all received recommendations")
    public ResponseEntity<?> listReceivedRecommendations(
            Pageable pageable,
            @Valid RecommendationReceivedRequest recommendationReceivedRequest,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(recommendationService.listReceivedRecommendations(
                    pageable,
                    recommendationReceivedRequest,
                    principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/sent")
    @Operation(summary = "Search among all sent recommendations")
    public ResponseEntity<?> listSentRecommendations(
            Pageable pageable,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(recommendationService.listSentRecommendations(pageable, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/")
    @Operation(summary = "Send new recommendation")
    public ResponseEntity<?> createRecommendation(
            @RequestBody @Valid RecommendationCreateRequest recommendationCreateRequest,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(recommendationService.createRecommendation(
                    recommendationCreateRequest,
                    principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
