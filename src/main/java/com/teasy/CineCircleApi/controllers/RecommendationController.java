package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationReceivedRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.RecommendationService;
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
public class RecommendationController {

    RecommendationService recommendationService;

    @Autowired
    public RecommendationController(RecommendationService recommendationService) {
        this.recommendationService = recommendationService;
    }

    @GetMapping("/received")
    public ResponseEntity<?> listReceivedRecommendations(
            Pageable pageable,
            RecommendationReceivedRequest recommendationReceivedRequest,
            Principal principal) {
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
    public ResponseEntity<?> listSentRecommendations(Pageable pageable, Principal principal) {
        try {
            return ResponseEntity.ok().body(recommendationService.listSentRecommendations(pageable, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/")
    public ResponseEntity<?> createRecommendation(
            @RequestBody RecommendationCreateRequest recommendationCreateRequest,
            Principal principal) {
        try {
            return ResponseEntity.ok().body(recommendationService.createRecommendation(
                    recommendationCreateRequest,
                    principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
