package com.teasy.CineCircleApi.controllers;

import com.teasy.CineCircleApi.models.dtos.requests.AuthMeUpdateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.AuthResetPasswordRequest;
import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.models.dtos.responses.AuthSignInResponse;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.TokenService;
import com.teasy.CineCircleApi.services.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/auth")
@CrossOrigin
@Tag(name = "Authentication", description = "All operations related to authenticated user")
public class AuthController {
    UserService userService;
    TokenService tokenService;

    @Autowired
    private AuthController(UserService userService,
                           TokenService tokenService) {
        this.userService = userService;
        this.tokenService = tokenService;
    }

    @GetMapping("/sign-in")
    @Operation(
            summary = "Get JWT token from user credentials",
            description = "Generate JWT token based on user credentials defined as Basic Auth in request header"
    )
    @SecurityRequirement(name = "basic")
    public ResponseEntity<?> createAuthenticationToken(Authentication authentication) {
        try {
            var jwtToken = tokenService.generateToken(authentication);
            var username = tokenService.getUsernameFromToken(jwtToken.tokenString());
            var user = userService.getUserFullInfoByUsernameOrEmail(username, username);
            return ResponseEntity.ok().body(new AuthSignInResponse(jwtToken, user));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/sign-up")
    @Operation(summary = "Create new user account")
    public ResponseEntity<?> createUser(
            @RequestBody @Valid AuthSignUpRequest request
    ) {
        try {
            return ResponseEntity.ok().body(userService.createUser(request));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/me")
    @Operation(summary = "Update authenticated user informations")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<?> updateUser(
            @RequestBody @Valid AuthMeUpdateRequest request,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(userService.updateUser(request, principal.getName()));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/me")
    @Operation(summary = "Get authenticated user informations")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<?> me(Principal principal) {
        try {
            return ResponseEntity.ok().body(userService.getUserFullInfo(principal.getName()));
        } catch (CustomException e) {
            // here we want to return 401 if user is not found
            if (e.getStatusCode() == HttpStatus.NOT_FOUND) {
                return new ResponseEntity<>(e.getMessage(), HttpStatus.UNAUTHORIZED);
            } else {
                return e.getEntityResponse();
            }
        }
    }

    @PostMapping("/reset-password")
    @Operation(summary = "Reset password for authenticated user (old password needed)")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<?> resetPassword(
            @RequestBody @Valid AuthResetPasswordRequest authResetPasswordRequest,
            Principal principal
    ) {
        try {
            return ResponseEntity.ok().body(userService.resetPassword(principal.getName(), authResetPasswordRequest));
        } catch (CustomException e) {
            // here we want to return 401 if user is not found
            if (e.getStatusCode() == HttpStatus.NOT_FOUND) {
                return new ResponseEntity<>(e.getMessage(), HttpStatus.UNAUTHORIZED);
            } else {
                return e.getEntityResponse();
            }
        }
    }
}
