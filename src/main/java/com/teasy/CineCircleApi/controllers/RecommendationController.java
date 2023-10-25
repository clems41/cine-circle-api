package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.CircleService;
import com.teasy.CineCircleApi.services.RecommendationService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

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
    public ResponseEntity<?> listReceivedRecommendations(Pageable pageable) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(recommendationService.listReceivedRecommendations(pageable, usernameFromToken));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/sent")
    public ResponseEntity<?> listSentRecommendations(Pageable pageable) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(recommendationService.listSentRecommendations(pageable, usernameFromToken));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/")
    public ResponseEntity<?> createRecommendation(@RequestBody RecommendationCreateRequest recommendationCreateRequest) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(recommendationService.createRecommendation(recommendationCreateRequest, usernameFromToken));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
