package com.teasy.CineCircleApi.controllers;

import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.AuthMeUpdateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.AuthRefreshTokenRequest;
import com.teasy.CineCircleApi.models.dtos.requests.AuthResetPasswordRequest;
import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.models.dtos.responses.AuthRefreshTokenResponse;
import com.teasy.CineCircleApi.models.dtos.responses.AuthSignInResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.TokenService;
import com.teasy.CineCircleApi.services.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
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
    public ResponseEntity<AuthSignInResponse> createAuthenticationToken(Authentication authentication) throws ExpectedException {
        var jwtToken = tokenService.generateToken(authentication.getName());
        var username = tokenService.getUsernameFromToken(jwtToken.tokenString());
        var user = userService.getUserFullInfoByUsernameOrEmail(username, username);
        var refreshToken = userService.getRefreshTokenForUser(username);
        return ResponseEntity.ok().body(new AuthSignInResponse(jwtToken, refreshToken, user));
    }

    @PostMapping("/sign-up")
    @Operation(summary = "Create new user account")
    public ResponseEntity<UserFullInfoDto> createUser(
            @RequestBody @Valid AuthSignUpRequest request
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.createUser(request));
    }

    @PostMapping("/refresh")
    @Operation(summary = "Refresh JWT token")
    public ResponseEntity<AuthRefreshTokenResponse> refreshToken(
            @RequestBody @Valid AuthRefreshTokenRequest request
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.refreshToken(request));
    }

    @PutMapping("/me")
    @Operation(summary = "Update authenticated user informations")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<UserFullInfoDto> updateUser(
            @RequestBody @Valid AuthMeUpdateRequest request,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.updateUser(request, principal.getName()));
    }

    @GetMapping("/me")
    @Operation(summary = "Get authenticated user informations")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<UserFullInfoDto> me(Principal principal) throws ExpectedException {
        return ResponseEntity.ok().body(userService.getUserFullInfo(principal.getName()));
    }

    @PostMapping("/reset-password")
    @Operation(summary = "Reset password for authenticated user (old password needed)")
    @SecurityRequirement(name = "JWT")
    public ResponseEntity<UserDto> resetPassword(
            @RequestBody @Valid AuthResetPasswordRequest authResetPasswordRequest,
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(userService.resetPassword(principal.getName(), authResetPasswordRequest));
    }
}
