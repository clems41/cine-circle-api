package com.teasy.CineCircleApi.controllers;

import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.services.HeadingService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.List;
import java.util.UUID;

@RestController
@Slf4j
@RequiredArgsConstructor
@RequestMapping("/api/v1/headings")
@CrossOrigin
@Tag(name = "Heading", description = "Add, remove and list headings from authenticated user, but also other users")
@SecurityRequirement(name = "JWT")
public class HeadingController {
    HeadingService headingService;

    @Autowired
    private HeadingController(HeadingService headingService) {
        this.headingService = headingService;
    }

    @GetMapping("/users/{userId}")
    @Operation(summary = "List headings for specific user")
    public ResponseEntity<List<MediaShortDto>> listHeadings(
            @PathVariable("userId") UUID userId
    ) throws ExpectedException {
        return ResponseEntity.ok().body(headingService.listHeadings(userId));
    }

    @GetMapping("/")
    @Operation(summary = "List headings for authenticated user")
    public ResponseEntity<List<MediaShortDto>> listHeadingsFOrAuthenticatedUser(
            Principal principal
    ) throws ExpectedException {
        return ResponseEntity.ok().body(headingService.listHeadingsForAuthenticatedUser(principal.getName()));
    }

    @PostMapping("/{mediaId}")
    @Operation(summary = "Add media to authenticated user headings")
    public ResponseEntity<String> addToLibrary(
            @PathVariable("mediaId") UUID mediaId,
            Principal principal
    ) throws ExpectedException {
        headingService.addToHeadings(mediaId, principal.getName());
        return ResponseEntity.ok().body("");
    }

    @DeleteMapping("/{mediaId}")
    @Operation(summary = "Remove media from authenticated user headings")
    public ResponseEntity<String> removeFromLibrary(
            @PathVariable("mediaId") UUID mediaId,
            Principal principal
    ) throws ExpectedException {
        headingService.removeFromHeadings(mediaId, principal.getName());
        return ResponseEntity.ok().body("");
    }
}
