package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
public class HeadingService {
    private final UserService userService;
    private final MediaService mediaService;
    private final static int MAX_HEADING_PER_USER = 3;

    @Autowired
    public HeadingService(UserService userService, MediaService mediaService) {
        this.mediaService = mediaService;
        this.userService = userService;
    }

    public void addToHeadings(UUID mediaId, String username) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);

        // if max has been already reached, delete the oldest one from headings list
        if (user.getHeadings().size() == MAX_HEADING_PER_USER) {
            var oldestOne = user.getHeadings().stream().toList().getLast();
            user.removeMediaFromHeadings(oldestOne);
        }

        // add media to headings
        var media = mediaService.findMediaByIdOrElseThrow(mediaId);
        user.addMediaToHeadings(media);
        userService.save(user);
    }

    public void removeFromHeadings(UUID mediaId, String username) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        var media = mediaService.findMediaByIdOrElseThrow(mediaId);
        user.removeMediaFromHeadings(media);
        userService.save(user);
    }

    public List<MediaShortDto> listHeadings(UUID userId) throws ExpectedException {
        var user = userService.findUserByIdOrElseThrow(userId);
        return user.getHeadings().stream().map(
                this::fromMediaEntityToMediaDto
        ).toList();
    }

    public List<MediaShortDto> listHeadingsForAuthenticatedUser(String username) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        return user.getHeadings().stream().map(
                this::fromMediaEntityToMediaDto
        ).toList();
    }

    private MediaShortDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaShortDto.class);
    }
}
