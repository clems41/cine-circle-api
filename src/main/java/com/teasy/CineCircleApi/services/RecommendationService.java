package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationSearchRequest;
import com.teasy.CineCircleApi.models.dtos.responses.RecommendationCreateResponse;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.RecommendationType;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.services.utils.CustomExampleMatcher;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.*;
import org.springframework.stereotype.Service;

import java.util.*;

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

    public RecommendationCreateResponse createRecommendation(RecommendationCreateRequest recommendationCreateRequest,
                                                             String authenticatedUsername) throws ExpectedException {
        var sentBy = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
        UUID recommendationRef = UUID.randomUUID(); // unique reference for all recommendations that will be created

        // adding receivers from list of users but also from list of circles
        Set<User> receivers = new HashSet<>();
        if (recommendationCreateRequest.userIds() != null) {
            for (UUID userId : recommendationCreateRequest.userIds()) {
                receivers.add(userService.findUserByIdOrElseThrow(userId));
            }
        }
        if (recommendationCreateRequest.circleIds() != null) {
            for (UUID circleId : recommendationCreateRequest.circleIds()) {
                var circle = circleService.findCircleByIdOrElseThrow(circleId);
                receivers.addAll(circle.getUsers());
            }
        }

        // find media
        var media = mediaService.findMediaByIdOrElseThrow(recommendationCreateRequest.mediaId());

        // create one recommendation for every receiver
        for (User user : receivers) {
            var recommendation = new Recommendation(
                    recommendationRef,
                    sentBy,
                    media,
                    user,
                    recommendationCreateRequest.comment(),
                    recommendationCreateRequest.rating()
            );
            recommendationRepository.save(recommendation);

            // send recommendation to concerned users
            notificationService.sendRecommendation(fromEntityToDto(recommendation));
        }

        // add media to library for sender
        libraryService.addToLibrary(media.getId(), null, authenticatedUsername);

        return new RecommendationCreateResponse(recommendationRef.toString());
    }

    public Page<RecommendationDto> searchRecommendations(
            Pageable pageable,
            RecommendationSearchRequest recommendationSearchRequest,
            String authenticatedUsername
    ) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);

        ExampleMatcher matcher = CustomExampleMatcher.matchingAll()
                .withIgnoreNullValues();
        var matchingRecommendation = new Recommendation();
        RecommendationType recommendationType;
        if (recommendationSearchRequest.type() == null) {
            recommendationType = RecommendationType.RECEIVED;
        } else {
            recommendationType = RecommendationType.getFromString(recommendationSearchRequest.type());
        }
        if (recommendationType == RecommendationType.RECEIVED) {
            matchingRecommendation.setReceiver(createMatchingUser(user));
        } else if (recommendationType == RecommendationType.SENT) {
            matchingRecommendation.setSentBy(createMatchingUser(user));
        } else {
            throw new ExpectedException(ErrorDetails.ERR_RECOMMENDATION_TYPE_NOT_SUPPORTED.addingArgs(recommendationSearchRequest.type()));
        }
        if (recommendationSearchRequest.mediaId() != null) {
            var media = mediaService.findMediaByIdOrElseThrow(recommendationSearchRequest.mediaId());
            matchingRecommendation.setMedia(createMatchingMedia(media));
        }
        if (recommendationSearchRequest.read() != null) {
            matchingRecommendation.setRead(recommendationSearchRequest.read());
        }
        return recommendationRepository.findAll(Example.of(matchingRecommendation, matcher), pageable).map(this::fromEntityToDto);
    }

    public void markRecommendationAsRead(UUID recommendationId, String authenticatedUsername) throws ExpectedException {
        // check that authenticated user is the receiver before marking the recommendation as read
        var recommendation = recommendationRepository
                .findById(recommendationId)
                .orElseThrow(() -> new ExpectedException(ErrorDetails.ERR_RECOMMENDATION_NOT_FOUND.addingArgs(recommendationId)));
        if (!recommendation.getReceiver().getUsername().equals(authenticatedUsername)) {
            throw new ExpectedException(ErrorDetails.ERR_RECOMMENDATION_RECEIVER_BAD_PERMISSIONS
                    .addingArgs(authenticatedUsername, recommendationId));
        }

        // mark as read
        recommendation.setRead(true);
        recommendationRepository.save(recommendation);
    }

    private User createMatchingUser(User user) {
        var matchingUser = new User();
        matchingUser.setId(user.getId());
        return matchingUser;
    }

    private Media createMatchingMedia(Media media) {
        var matchingMedia = new Media();
        matchingMedia.setId(media.getId());
        return matchingMedia;
    }

    private RecommendationDto fromEntityToDto(Recommendation recommendation) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(recommendation, RecommendationDto.class);
    }
}
