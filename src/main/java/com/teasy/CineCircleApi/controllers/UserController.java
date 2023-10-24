package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.UserResetPasswordRequest;
import com.teasy.CineCircleApi.models.dtos.requests.UserSearchRequest;
import com.teasy.CineCircleApi.models.dtos.requests.UserSendResetPasswordEmailRequest;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.services.UserService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/users")
@CrossOrigin
public class UserController {
    UserService userService;

    @Autowired
    private UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("/")
    public ResponseEntity<?> searchUsers(Pageable page, UserSearchRequest request) {
        try {
            return ResponseEntity.ok().body(userService.searchUsers(page, request));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/reset-password")
    public ResponseEntity<?> sendResetPasswordEmail(UserSendResetPasswordEmailRequest userSendResetPasswordEmailRequest) {
        try {
            userService.sendResetPasswordEmail(userSendResetPasswordEmailRequest.email());
            return ResponseEntity.ok().body("");
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @PostMapping("/reset-password")
    public ResponseEntity<?> resetPassword(@RequestBody UserResetPasswordRequest userResetPasswordRequest) {
        try {
            return ResponseEntity.ok().body(userService.resetPasswordWithToken(userResetPasswordRequest));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<?> getMedia(final @PathVariable Long id) {
        try {
            return ResponseEntity.ok().body(userService.getUser(id));
        } catch (CustomException e) {
            return e.getEntityResponse();
        }
    }
}
