package com.teasy.CineCircleApi.controllers;

import com.teasy.CineCircleApi.models.dtos.UserDetails;
import com.teasy.CineCircleApi.models.dtos.responses.SignInResponse;
import com.teasy.CineCircleApi.services.TokenService;
import com.teasy.CineCircleApi.services.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
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
        var token = tokenService.generateToken(authentication);
        return ResponseEntity.ok().body(new SignInResponse(token, new UserDetails()));
    }
}
