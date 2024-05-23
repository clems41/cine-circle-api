package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.NotificationService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
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
    public ResponseEntity<NotificationTopicResponse> getTopicForUser(Principal principal) throws ExpectedException {
        return ResponseEntity.ok().body(notificationService.getTopicForUsername(principal.getName()));
    }
}
