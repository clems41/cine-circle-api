package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.dtos.responses.RecommendationCreateResponse;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.utils.CustomPageImpl;
import com.teasy.CineCircleApi.utils.HttpUtils;
import com.teasy.CineCircleApi.utils.RandomUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.http.*;

import java.time.LocalDateTime;
import java.util.*;
import java.util.stream.Stream;

public class RecommendationTest extends IntegrationTestAbstract {
    @Test
    public void sendRecommendationToMultipleUsers() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var receiver1 = dummyDataCreator.generateUser(true);
        var receiver2 = dummyDataCreator.generateUser(true);
        var receiver3 = dummyDataCreator.generateUser(true);
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var comment = RandomUtils.randomString(50);
        var rating = RandomUtils.randomInt(1, 5);

        /* Send recommendation */
        var userIds = List.of(receiver1.getId(), receiver2.getId(), receiver3.getId());
        var recommendationCreateRequest = new RecommendationCreateRequest(
                media.getId(), userIds, List.of(), comment, rating);
        ResponseEntity<RecommendationCreateResponse> recommendationCreateResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(recommendationCreateRequest, headers),
                        RecommendationCreateResponse.class
                );
        // check response
        Assertions.assertThat(recommendationCreateResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(recommendationCreateResponse.getBody()).isNotNull();
        var recommendationRef = recommendationCreateResponse.getBody().recommendationRef();

        /* Check data in database that all recommendations have been created */
        var recommendations = recommendationRepository
                .findAll()
                .stream()
                .filter(recommendation -> recommendation.getSentBy().getId().equals(authenticatedUser.getId()))
                .toList();
        Assertions.assertThat(recommendations).hasSize(3);
        List<UUID> receiverIdsFromDatabase = new ArrayList<>();
        for (var recommendation : recommendations) {
            Assertions.assertThat(recommendation.getRecommendationRef().toString()).isEqualTo(recommendationRef);
            Assertions.assertThat(recommendation.getComment()).isEqualTo(comment);
            Assertions.assertThat(recommendation.getRating()).isEqualTo(rating);
            Assertions.assertThat(recommendation.getRead()).isEqualTo(false);
            Assertions.assertThat(recommendation.getSentBy().getId()).isEqualTo(authenticatedUser.getId());
            Assertions.assertThat(recommendation.getMedia().getId()).isEqualTo(media.getId());
            Assertions.assertThat(recommendation.getSentAt()).isBefore(LocalDateTime.now().plusMinutes(1));
            Assertions.assertThat(recommendation.getSentAt()).isAfter(LocalDateTime.now().minusMinutes(1));
            receiverIdsFromDatabase.add(recommendation.getReceiver().getId());
        }
        Assertions.assertThat(receiverIdsFromDatabase).containsExactlyInAnyOrderElementsOf(userIds);
    }

    @Test
    public void markRecommendationAsRead() {
        /* Data */
        var signUpRequestForSender = authenticator.authenticateNewUser();
        var signUpRequestForReceiver = authenticator.authenticateNewUser();
        var headersForReceiver = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequestForReceiver.username(), signUpRequestForReceiver.password());
        var sender = userRepository.findByUsername(signUpRequestForSender.username()).orElseThrow();
        var receiver = userRepository.findByUsername(signUpRequestForReceiver.username()).orElseThrow();
        var otherReceiver1 = dummyDataCreator.generateUser(true);
        var otherReceiver2 = dummyDataCreator.generateUser(true);
        var otherReceiver3 = dummyDataCreator.generateUser(true);
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var comment = RandomUtils.randomString(50);
        var rating = RandomUtils.randomInt(1, 5);
        var recommendationRef = UUID.randomUUID();

        /* Create recommendations in database */
        var recommendationForReceiver = new Recommendation(recommendationRef, sender, media, receiver, comment, rating);
        recommendationRepository.save(recommendationForReceiver);
        recommendationRepository.save(new Recommendation(recommendationRef, sender, media, otherReceiver1, comment, rating));
        recommendationRepository.save(new Recommendation(recommendationRef, sender, media, otherReceiver2, comment, rating));
        recommendationRepository.save(new Recommendation(recommendationRef, sender, media, otherReceiver3, comment, rating));

        /* Mark recommendation as read for one receiver */
        ResponseEntity<String> markAsReadResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.recommendationUrl)
                                .concat(String.format("/%s/read", recommendationForReceiver.getId())),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headersForReceiver),
                        String.class
                );
        Assertions.assertThat(markAsReadResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Check data in database that the recommendation for receiver has been marked as read */
        var recommendation = recommendationRepository.findById(recommendationForReceiver.getId()).orElseThrow();
            Assertions.assertThat(recommendation.getRecommendationRef()).isEqualTo(recommendationRef);
            Assertions.assertThat(recommendation.getComment()).isEqualTo(comment);
            Assertions.assertThat(recommendation.getRating()).isEqualTo(rating);
            Assertions.assertThat(recommendation.getSentBy().getId()).isEqualTo(sender.getId());
            Assertions.assertThat(recommendation.getMedia().getId()).isEqualTo(media.getId());
            Assertions.assertThat(recommendation.getSentAt()).isBefore(LocalDateTime.now().plusMinutes(1));
            Assertions.assertThat(recommendation.getSentAt()).isAfter(LocalDateTime.now().minusMinutes(1));
            Assertions.assertThat(recommendation.getRead()).isEqualTo(true);
    }

    @ParameterizedTest
    @MethodSource("searchRecommendations_args")
    public void searchRecommendations(String recommendationType, Boolean read, Pageable pageable,
                                      int nbExpectedForAll, int nbExpectedForSpecificMedia) {
        /* Create data in database */
        var specificMedia = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var headersToUse = generateRecommendationsAndGetHttpHeadersForAuthenticatedUser(specificMedia);

        /* Send request with query parameters without specifying media */
        Map<String, Object> queryParams = new HashMap<>();
        if (pageable != null) {
            queryParams.put("page", pageable.getPageNumber());
            queryParams.put("size", pageable.getPageSize());
        }
        if (recommendationType != null) {
            queryParams.put("type", recommendationType);
        }
        if (read != null) {
            queryParams.put("read", read);
        }
        ResponseEntity<CustomPageImpl<RecommendationDto>> responseForAll = this.restTemplate
                .exchange(HttpUtils.getUriWithQueryParameter(port, HttpUtils.recommendationUrl, queryParams),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headersToUse),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(responseForAll.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(responseForAll.getBody()).isNotNull();
        Assertions.assertThat(responseForAll.getBody().stream().toList()).hasSize(nbExpectedForAll);

        /* Send request with query parameters with specific media */
        queryParams.put("mediaId", specificMedia.getId());
        ResponseEntity<CustomPageImpl<RecommendationDto>> responseForSpecificMedia = this.restTemplate
                .exchange(HttpUtils.getUriWithQueryParameter(port, HttpUtils.recommendationUrl, queryParams),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headersToUse),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(responseForSpecificMedia.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(responseForSpecificMedia.getBody()).isNotNull();
        Assertions.assertThat(responseForSpecificMedia.getBody().stream().toList()).hasSize(nbExpectedForSpecificMedia);
    }


        /*** List of expectations :
        - 18 received in total : 10 seen for any media + 3 seen for one specific media + 4 unseen for any media + 1 unseen for specific media
        - 12 sent in total : 4 seen for any media + 5 seen for one specific media + 3 unseen for any media + 0 unseen for specific media
        ***/

    static Stream<Arguments> searchRecommendations_args() {
        String receivedRecommendationType = "received";
        String sentRecommendationType = "sent";
        var pageable = PageRequest.of(0, 20);
        return Stream.of(
                Arguments.of(null, null, null, 10, 4), // should be 18, but the default pageSize is 10, so the max expected is 10
                Arguments.of(receivedRecommendationType, null, pageable, 18, 4),
                Arguments.of(sentRecommendationType, null, pageable, 12, 5),
                Arguments.of(receivedRecommendationType, true, pageable, 13, 3),
                Arguments.of(sentRecommendationType, false, pageable, 3, 0)
        );
    }

    private HttpHeaders generateRecommendationsAndGetHttpHeadersForAuthenticatedUser(Media media) {
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var user = userRepository.findByUsername(signUpRequest.username()).orElseThrow();

        // random recommendations unseen
        for (int i = 0; i < 20; i++) {
            dummyDataCreator.generateRecommendation(true, null, null, null, false);
        }
        // random recommendations seen
        for (int i = 0; i < 20; i++) {
            dummyDataCreator.generateRecommendation(true, null, null, null, true);
        }

        // recommendations received seen by the authenticated user
        for (int i = 0; i < 10; i++) {
            dummyDataCreator.generateRecommendation(true, null, user, null, true);
        }

        // recommendations received seen for specific mediaId by the authenticated user
        for (int i = 0; i < 3; i++) {
            dummyDataCreator.generateRecommendation(true, null, user, media, true);
        }

        // recommendations received unseen by the authenticated user
        for (int i = 0; i < 4; i++) {
            dummyDataCreator.generateRecommendation(true, null, user, null, false);
        }

        // recommendations received unseen for specific mediaId by the authenticated user
        for (int i = 0; i < 1; i++) {
            dummyDataCreator.generateRecommendation(true, null, user, media, false);
        }

        // recommendations sent seen by the authenticated user
        for (int i = 0; i < 4; i++) {
            dummyDataCreator.generateRecommendation(true, user, null, null, true);
        }

        // recommendations sent seen for specific mediaId by the authenticated user
        for (int i = 0; i < 5; i++) {
            dummyDataCreator.generateRecommendation(true, user, null, media, true);
        }

        // recommendations sent unseen by the authenticated user
        for (int i = 0; i < 3; i++) {
            dummyDataCreator.generateRecommendation(true, user, null, null, false);
        }

        // recommendations sent unseen for specific mediaId by the authenticated user
        for (int i = 0; i < 0; i++) {
            dummyDataCreator.generateRecommendation(true, user, null, media, false);
        }

        return headers;
    }
}
