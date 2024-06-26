package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationReceivedRequest;
import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.ErrorMessage;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.CircleRepository;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Service
@Slf4j
public class RecommendationService {
    private final NotificationServiceInterface notificationService;
    private final RecommendationRepository recommendationRepository;
    private final UserRepository userRepository;
    private final MediaRepository mediaRepository;
    private final LibraryService libraryService;
    private final CircleRepository circleRepository;

    @Autowired
    public RecommendationService(RecommendationRepository recommendationRepository,
                                 UserRepository userRepository,
                                 MediaRepository mediaRepository,
                                 LibraryService libraryService,
                                 CircleRepository circleRepository,
                                 NotificationServiceInterface notificationService) {
        this.recommendationRepository = recommendationRepository;
        this.userRepository = userRepository;
        this.mediaRepository = mediaRepository;
        this.notificationService = notificationService;
        this.libraryService = libraryService;
        this.circleRepository = circleRepository;
    }

    public RecommendationDto createRecommendation(RecommendationCreateRequest recommendationCreateRequest,
                                                  String authenticatedUsername) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(authenticatedUsername);

        // adding receivers
        Set<User> receivers = new HashSet<>();
        if (recommendationCreateRequest.userIds() != null) {
            recommendationCreateRequest.userIds()
                    .forEach(userId -> {
                        try {
                            receivers.add(findUserByIdOrElseThrow(userId));
                        } catch (ExpectedException e) {
                            log.error("User not found", e);
                        }
                    });
        }

        // adding circles
        Set<Circle> circles = new HashSet<>();
        if (recommendationCreateRequest.circleIds() != null) {
            recommendationCreateRequest.circleIds()
                    .forEach(circleId -> {
                        try {
                            circles.add(findCircleByIdOrElseThrow(circleId));
                        } catch (ExpectedException e) {
                            throw new RuntimeException(e);
                        }
                    });
        }

        // find media
        var media = findMediaByIdOrElseThrow(recommendationCreateRequest.mediaId());

        // create recommendation
        var recommendation = new Recommendation(
                user,
                media,
                receivers,
                circles,
                recommendationCreateRequest.comment(),
                recommendationCreateRequest.rating());
        recommendationRepository.save(recommendation);

        // add media to library for sender
        libraryService.addToLibrary(media.getId(), null, authenticatedUsername);

        // send recommendation to concerned users
        var recommendationDto = fromEntityToDto(recommendation);
        notificationService.sendRecommendation(recommendationDto);

        return recommendationDto;
    }

    public Page<RecommendationDto> listReceivedRecommendations(
            Pageable pageable,
            RecommendationReceivedRequest recommendationReceivedRequest,
            String authenticatedUsername) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(authenticatedUsername);
        Page<Recommendation> result;
        if (recommendationReceivedRequest.mediaId() != null) {
            var media = findMediaByIdOrElseThrow(recommendationReceivedRequest.mediaId());
            result = recommendationRepository.findAllByReceivers_IdAndMedia_Id(pageable, user.getId(), media.getId());
        } else {
            result = recommendationRepository.findAllByReceivers_Id(pageable, user.getId());
        }
        return result.map(this::fromEntityToDto);
    }

    public Page<RecommendationDto> listSentRecommendations(Pageable pageable, String authenticatedUsername) throws ExpectedException {
        // creating matching example
        var user = findUserByUsernameOrElseThrow(authenticatedUsername);
        var matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        User matchingSender = new User();
        matchingSender.setId(user.getId());
        var matchingRecommendation = new Recommendation(matchingSender, null, null, null, null, null);
        matchingRecommendation.setSentAt(null);

        var result = recommendationRepository.findAll(Example.of(matchingRecommendation, matcher), pageable);
        return result.map(this::fromEntityToDto);
    }

    private RecommendationDto fromEntityToDto(Recommendation recommendation) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(recommendation, RecommendationDto.class);
    }

    private Circle findCircleByIdOrElseThrow(UUID circleId) throws ExpectedException {
        // check if media exists
        return circleRepository
                .findById(circleId)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.CIRCLE_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private Media findMediaByIdOrElseThrow(UUID mediaId) throws ExpectedException {
        // check if media exists
        return mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.MEDIA_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private User findUserByIdOrElseThrow(UUID id) throws ExpectedException {
        return userRepository
                .findById(id)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.USER_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private User findUserByUsernameOrElseThrow(String username) throws ExpectedException {
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.USER_NOT_FOUND, HttpStatus.NOT_FOUND));
    }
}
