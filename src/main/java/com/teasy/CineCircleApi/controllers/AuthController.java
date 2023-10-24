package com.teasy.CineCircleApi.controllers;

import com.teasy.CineCircleApi.models.dtos.requests.AuthMeUpdateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.AuthResetPasswordRequest;
import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.models.dtos.responses.SignInResponse;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.TokenService;
import com.teasy.CineCircleApi.services.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/auth")
@CrossOrigin
public class AuthController {
    UserService userService;
    TokenService tokenService;

    @Autowired
    private AuthController(UserService userService,
                           TokenService tokenService) {
        this.userService = userService;
        this.tokenService = tokenService;
    }

    @PostMapping("/sign-in")
    public ResponseEntity<?> createAuthenticationToken(Authentication authentication) {
        try {
            var jwtToken = tokenService.generateToken(authentication);
            var username = tokenService.getUsernameFromToken(jwtToken.tokenString());
            var user = userService.getUserFullInfoByUsernameOrEmail(username, username);
            return ResponseEntity.ok().body(new SignInResponse(jwtToken, user));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/sign-up")
    public ResponseEntity<?> createUser(@RequestBody AuthSignUpRequest request) {
        try {
            return ResponseEntity.ok().body(userService.createUser(request));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PutMapping("/me")
    public ResponseEntity<?> updateUser(@RequestBody AuthMeUpdateRequest request) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(userService.updateUser(request, usernameFromToken));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/me")
    public ResponseEntity<?> me() {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(userService.getUserFullInfo(usernameFromToken));
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
    public ResponseEntity<?> resetPassword(@RequestBody AuthResetPasswordRequest authResetPasswordRequest) {
        try {
            var usernameFromToken = SecurityContextHolder
                    .getContext()
                    .getAuthentication()
                    .getName();
            return ResponseEntity.ok().body(userService.resetPassword(usernameFromToken, authResetPasswordRequest));
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
