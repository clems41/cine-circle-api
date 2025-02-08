package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.UserResetPasswordRequest;
import com.teasy.CineCircleApi.models.dtos.requests.UserSearchRequest;
import com.teasy.CineCircleApi.models.dtos.requests.UserSendResetPasswordEmailRequest;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.UserService;
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
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/users")
@CrossOrigin
@Tag(name = "User", description = "Search among users or reset password (no old password needed)")
public class UserController {
    UserService userService;

    @Autowired
    private UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("")
    @Operation(summary = "Search for user")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<Page<UserDto>> searchUsers(
            Pageable page,
            @Valid UserSearchRequest request,
            Authentication authentication
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.searchUsers(authentication.getName(), page, request));
    }

    @GetMapping("/reset-password")
    @Operation(summary = "Send email with unique token to reset password for an unauthenticated user")
    public ResponseEntity<String> sendResetPasswordEmail(
            @Valid UserSendResetPasswordEmailRequest userSendResetPasswordEmailRequest
    ) throws ExpectedException {
        userService.sendResetPasswordEmail(userSendResetPasswordEmailRequest.email());
        return ResponseEntity.ok().body("");
    }

    @PostMapping("/reset-password")
    @Operation(summary = "Reset password with unique token received by email for an unauthenticated user")
    public ResponseEntity<UserDto> resetPassword(
            @RequestBody @Valid UserResetPasswordRequest userResetPasswordRequest
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.resetPasswordWithToken(userResetPasswordRequest));
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get details about specific user")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<UserDto> getUser(
            @PathVariable UUID id
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.getUser(id));
    }

    @PutMapping("/related/{related_user_id}")
    @Operation(summary = "Add user in related users")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<UserFullInfoDto> addUserInRelatedUsers(
            Principal principal,
            @PathVariable UUID related_user_id
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.addUserToRelatedUsers(principal.getName(), related_user_id));
    }

    @DeleteMapping("/related/{related_user_id}")
    @Operation(summary = "Remove user from related users")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<UserFullInfoDto> removeUserFromRelatedUsers(
            Principal principal,
            @PathVariable UUID related_user_id
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.removeUserFromRelatedUsers(principal.getName(), related_user_id));
    }
}
