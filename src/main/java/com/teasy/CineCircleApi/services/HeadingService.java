package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.ErrorMessage;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import java.util.Comparator;
import java.util.List;
import java.util.UUID;

@Service
public class HeadingService {
    UserRepository userRepository;
    MediaRepository mediaRepository;
    private final static int MAX_HEADING_PER_USER = 3;

    @Autowired
    public HeadingService(MediaRepository mediaRepository, UserRepository userRepository) {
        this.mediaRepository = mediaRepository;
        this.userRepository = userRepository;
    }

    public void addToHeadings(UUID mediaId, String username) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);

        // if max has been already reached, delete the oldest one from headings list
        if (user.getHeadings().size() == MAX_HEADING_PER_USER) {
            var oldestOne = user.getHeadings().stream().toList().getLast();
            user.removeMediaFromHeadings(oldestOne);
        }

        // add media to headings
        var media = findMediaByIdOrElseThrow(mediaId);
        user.addMediaToHeadings(media);
        userRepository.save(user);
    }

    public void removeFromHeadings(UUID mediaId, String username) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        var media = findMediaByIdOrElseThrow(mediaId);
        user.removeMediaFromHeadings(media);
        userRepository.save(user);
    }

    public List<MediaShortDto> listHeadings(UUID userId) throws ExpectedException {
        var user = findUserByIdOrElseThrow(userId);
        return user.getHeadings().stream().map(
                this::fromMediaEntityToMediaDto
        ).toList();
    }

    public List<MediaShortDto> listHeadingsForAuthenticatedUser(String username) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        return user.getHeadings().stream().map(
                this::fromMediaEntityToMediaDto
        ).toList();
    }

    private User findUserByUsernameOrElseThrow(String username) throws ExpectedException {
        // check if user exist
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.USER_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private User findUserByIdOrElseThrow(UUID userId) throws ExpectedException {
        // check if user exist
        return userRepository
                .findById(userId)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.USER_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private Media findMediaByIdOrElseThrow(UUID mediaId) throws ExpectedException {
        // check if media exists
        return mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.MEDIA_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private MediaShortDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaShortDto.class);
    }
}
