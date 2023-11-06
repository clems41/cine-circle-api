package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.UserResetPasswordRequest;
import com.teasy.CineCircleApi.models.dtos.requests.UserSearchRequest;
import com.teasy.CineCircleApi.models.dtos.requests.UserSendResetPasswordEmailRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/users")
@CrossOrigin
@Tag(name = "Library", description = "Search among users or reset password (no old password needed)")
public class UserController {
    UserService userService;

    @Autowired
    private UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("/")
    @Operation(summary = "Search for user")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<?> searchUsers(Pageable page, UserSearchRequest request) {
        try {
            return ResponseEntity.ok().body(userService.searchUsers(page, request));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/reset-password")
    @Operation(summary = "Send email with unique token to reset password for an unauthenticated user")
    public ResponseEntity<?> sendResetPasswordEmail(UserSendResetPasswordEmailRequest userSendResetPasswordEmailRequest) {
        try {
            userService.sendResetPasswordEmail(userSendResetPasswordEmailRequest.email());
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/reset-password")
    @Operation(summary = "Reset password with unique token received by email for an unauthenticated user")
    public ResponseEntity<?> resetPassword(@RequestBody UserResetPasswordRequest userResetPasswordRequest) {
        try {
            return ResponseEntity.ok().body(userService.resetPasswordWithToken(userResetPasswordRequest));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get details about specific user")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<?> getUser(final @PathVariable UUID id) {
        try {
            return ResponseEntity.ok().body(userService.getUser(id));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
