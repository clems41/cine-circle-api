package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationReceivedRequest;
import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.NotificationService;
import com.teasy.CineCircleApi.services.RecommendationService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.messaging.handler.annotation.MessageMapping;
import org.springframework.messaging.simp.annotation.SendToUser;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;

@RestController
@Slf4j
@RequiredArgsConstructor
@CrossOrigin
@RequestMapping("api/v1/notifications")
@Tag(name = "Notification", description = "Get notification topic to create Web Socket connection")
@SecurityRequirement(name = "JWT")
public class NotificationController {

    NotificationService notificationService;

    @Autowired
    public NotificationController(NotificationService notificationService) {
        this.notificationService = notificationService;
    }

    @GetMapping("/topic")
    @Operation(summary = "Get unique topic that should be used by authenticated user to get his notifications")
    public ResponseEntity<?> getTopicForUser(Principal principal) {
        try {
            return ResponseEntity.ok().body(notificationService.getTopicForUsername(principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
