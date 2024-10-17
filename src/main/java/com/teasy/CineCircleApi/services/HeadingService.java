package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.entities.Heading;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.ErrorMessage;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.HeadingRepository;
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
    HeadingRepository headingRepository;
    MediaRepository mediaRepository;
    UserRepository userRepository;
    private final static int MAX_HEADING_PER_USER = 3;

    @Autowired
    public HeadingService(HeadingRepository headingRepository, MediaRepository mediaRepository, UserRepository userRepository) {
        this.headingRepository = headingRepository;
        this.mediaRepository = mediaRepository;
        this.userRepository = userRepository;
    }

    public void addToHeadings(UUID mediaId, String username) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        var headings = headingRepository.findAllByUserId(user.getId());

        // if max has been already reached, delete the oldest one from headings list
        if (headings.size() == MAX_HEADING_PER_USER) {
            headings.sort(Comparator.comparing(Heading::getAddedAt));
            var oldestOne = headings.getLast();
            removeFromHeadings(oldestOne.getId(), username);
        }

        headingRepository.save(newHeading(username, mediaId));
    }

    public void removeFromHeadings(UUID mediaId, String username) throws ExpectedException {
        headingRepository.delete(newHeading(username, mediaId));
    }

    public List<MediaShortDto> listHeadings(UUID userId) throws ExpectedException {
        findUserByIdOrElseThrow(userId);
        var headings = headingRepository.findAllByUserId(userId);
        return headings.stream().map(
                heading -> fromMediaEntityToMediaDto(heading.getMedia())
        ).toList();
    }

    public List<MediaShortDto> listHeadingsForAuthenticatedUser(String username) throws ExpectedException {
        return listHeadings(findUserByUsernameOrElseThrow(username).getId());
    }

    private Heading newHeading(String username, UUID mediaId) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        var media = findMediaByIdOrElseThrow(mediaId);
        return new Heading(user, media);
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
