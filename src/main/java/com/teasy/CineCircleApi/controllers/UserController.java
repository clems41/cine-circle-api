package com.teasy.CineCircleApi.controllers;


import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.dtos.requests.UserSearchRequest;
import com.teasy.CineCircleApi.services.HttpErrorService;
import com.teasy.CineCircleApi.services.UserService;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.externals.mediaProviders.theMovieDb.TheMovieDb;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

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
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<?> getMedia(final @PathVariable Long id) {
        try {
            return ResponseEntity.ok().body(userService.getUser(id));
        } catch (ResponseStatusException e) {
            return HttpErrorService.getEntityResponseFromException(e);
        }
    }
}
