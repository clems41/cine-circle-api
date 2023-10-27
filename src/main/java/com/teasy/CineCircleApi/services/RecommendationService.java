package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationReceivedRequest;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.HashSet;
import java.util.Set;

@Service
public class RecommendationService {
    private RecommendationRepository recommendationRepository;
    private UserRepository userRepository;
    private MediaRepository mediaRepository;

    @Autowired
    public RecommendationService(RecommendationRepository recommendationRepository,
                                 UserRepository userRepository,
                                 MediaRepository mediaRepository) {
        this.recommendationRepository = recommendationRepository;
        this.userRepository = userRepository;
        this.mediaRepository = mediaRepository;
    }

    public RecommendationDto createRecommendation(RecommendationCreateRequest recommendationCreateRequest,
                                                  String authenticatedUsername) {
        var user = getUserWithUsernameOrElseThrow(authenticatedUsername);

        // adding receivers
        Set<User> receivers = new HashSet<>();
        recommendationCreateRequest.userIds()
                .forEach(userId -> receivers.add(getUserWithIdOrElseThrow(userId)));

        // find media
        var media = getMediaWithIdOrElseThrow(recommendationCreateRequest.mediaId());

        // create recommendation
        var recommendation = new Recommendation(
                user,
                media,
                receivers,
                recommendationCreateRequest.comment(),
                recommendationCreateRequest.rating());
        recommendationRepository.save(recommendation);
        return fromEntityToDto(recommendation);
    }

    public Page<RecommendationDto> listReceivedRecommendations(
            Pageable pageable,
            RecommendationReceivedRequest recommendationReceivedRequest,
            String authenticatedUsername) {
        var user = getUserWithUsernameOrElseThrow(authenticatedUsername);
        Page<Recommendation> result;
        if (recommendationReceivedRequest.mediaId() != null) {
            var media = getMediaWithIdOrElseThrow(recommendationReceivedRequest.mediaId());
            result = recommendationRepository.findAllByReceivers_IdAndMedia_Id(pageable, user.getId(), media.getId());
        } else {
            result = recommendationRepository.findAllByReceivers_Id(pageable, user.getId());
        }
        return result.map(this::fromEntityToDto);
    }

    public Page<RecommendationDto> listSentRecommendations(Pageable pageable, String authenticatedUsername) {
        // creating matching example
        var user = getUserWithUsernameOrElseThrow(authenticatedUsername);
        var matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        User matchingSender = new User();
        matchingSender.setId(user.getId());
        var matchingRecommendation = new Recommendation(matchingSender, null, null, null, null);
        matchingRecommendation.setSentAt(null);

        var result = recommendationRepository.findAll(Example.of(matchingRecommendation, matcher), pageable);
        return result.map(this::fromEntityToDto);
    }

    public RecommendationDto fromEntityToDto(Recommendation recommendation) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(recommendation, RecommendationDto.class);
    }

    private Media getMediaWithIdOrElseThrow(Long mediaId) {
        // check if media exists
        return mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> CustomExceptionHandler.mediaWithIdNotFound(mediaId));
    }

    private User getUserWithIdOrElseThrow(Long id) throws CustomException {
        return userRepository
                .findById(id)
                .orElseThrow(() -> CustomExceptionHandler.userWithIdNotFound(id));
    }

    private User getUserWithUsernameOrElseThrow(String username) throws CustomException {
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
    }
}
