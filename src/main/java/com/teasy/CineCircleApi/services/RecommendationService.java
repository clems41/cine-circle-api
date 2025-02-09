package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationReceivedRequest;
import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.*;
import org.springframework.stereotype.Service;

import java.util.HashSet;
import java.util.Set;

@Service
@Slf4j
public class RecommendationService {
    private final NotificationService notificationService;
    private final RecommendationRepository recommendationRepository;
    private final UserService userService;
    private final MediaService mediaService;
    private final LibraryService libraryService;
    private final CircleService circleService;

    @Autowired
    public RecommendationService(RecommendationRepository recommendationRepository,
                                 UserService userService,
                                 MediaService mediaService,
                                 LibraryService libraryService,
                                 CircleService circleService,
                                 NotificationService notificationService) {
        this.recommendationRepository = recommendationRepository;
        this.userService = userService;
        this.mediaService = mediaService;
        this.notificationService = notificationService;
        this.libraryService = libraryService;
        this.circleService = circleService;
    }

    public RecommendationDto createRecommendation(RecommendationCreateRequest recommendationCreateRequest,
                                                  String authenticatedUsername) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);

        // adding receivers
        Set<User> receivers = new HashSet<>();
        if (recommendationCreateRequest.userIds() != null) {
            recommendationCreateRequest.userIds()
                    .forEach(userId -> {
                        try {
                            receivers.add(userService.findUserByIdOrElseThrow(userId));
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
                            circles.add(circleService.findCircleByIdOrElseThrow(circleId));
                        } catch (ExpectedException e) {
                            throw new RuntimeException(e);
                        }
                    });
        }

        // find media
        var media = mediaService.findMediaByIdOrElseThrow(recommendationCreateRequest.mediaId());

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
        notificationService.sendRecommendation(recommendation);

        return fromEntityToDto(recommendation);
    }

    public Page<RecommendationDto> listReceivedRecommendations(
            Pageable pageable,
            RecommendationReceivedRequest recommendationReceivedRequest,
            String authenticatedUsername) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
        Page<Recommendation> result;
        if (recommendationReceivedRequest.mediaId() != null) {
            var media = mediaService.findMediaByIdOrElseThrow(recommendationReceivedRequest.mediaId());
            result = recommendationRepository.findAllByReceivers_IdAndMedia_Id(pageable, user.getId(), media.getId());
        } else {
            result = recommendationRepository.findAllByReceivers_Id(pageable, user.getId());
        }
        return result.map(this::fromEntityToDto);
    }

    public Page<RecommendationDto> listSentRecommendations(Pageable pageable, String authenticatedUsername) throws ExpectedException {
        // creating matching example
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
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
}
