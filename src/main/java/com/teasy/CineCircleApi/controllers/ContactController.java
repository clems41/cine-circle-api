package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.ContactSendFeedbackRequest;
import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.ContactService;
import com.teasy.CineCircleApi.services.NotificationService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;

@RestController
@Slf4j
@RequiredArgsConstructor
@CrossOrigin
@RequestMapping("api/v1/contact")
@Tag(name = "Contact", description = "Send feedback to the dev team")
@SecurityRequirement(name = "JWT")
public class ContactController {

    ContactService contactService;

    @Autowired
    public ContactController(ContactService contactService) {
        this.contactService = contactService;
    }

    @PostMapping("")
    @Operation(summary = "Send feedback")
    public ResponseEntity<String> sendFeedback(
            Principal principal,
            @Valid @RequestBody ContactSendFeedbackRequest request
    ) throws ExpectedException {
        contactService.sendFeedback(request, principal.getName());
        return ResponseEntity.ok().body("");
    }
}
