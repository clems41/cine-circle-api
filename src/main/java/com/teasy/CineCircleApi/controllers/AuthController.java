package com.teasy.CineCircleApi.controllers;

import com.teasy.CineCircleApi.models.entities.UserDetails;
import com.teasy.CineCircleApi.models.dtos.requests.SignUpRequest;
import com.teasy.CineCircleApi.models.dtos.responses.SignInResponse;
import com.teasy.CineCircleApi.services.HttpErrorService;
import com.teasy.CineCircleApi.services.TokenService;
import com.teasy.CineCircleApi.services.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

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
        var token = tokenService.generateToken(authentication);
        var username = tokenService.getUsernameFromToken(token);
        var user = userService.getUserByUsernameOrEmail(username, username);
        return ResponseEntity.ok().body(new SignInResponse(token, user));
    }

    @PostMapping("/sign-up")
    public ResponseEntity<?> createUser(@RequestBody SignUpRequest request) {
        try {
            return ResponseEntity.ok().body(userService.createUser(request));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @PostMapping("/hello")
    public ResponseEntity<?> hello() {
        return ResponseEntity.ok().body(
                String.format("hello authenticated user %s",
                        SecurityContextHolder
                        .getContext()
                        .getAuthentication()
                        .getName()));
    }
}
