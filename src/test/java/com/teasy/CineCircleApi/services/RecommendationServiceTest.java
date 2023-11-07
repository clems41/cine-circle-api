package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.config.RsaKeyProperties;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.mocks.NotificationServiceMock;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.context.ActiveProfiles;

import java.util.ArrayList;
import java.util.List;

@DataJpaTest
@ActiveProfiles("test")
public class RecommendationServiceTest {
    @MockBean
    RsaKeyProperties rsaKeyProperties;
    NotificationServiceMock notificationService;
    @Autowired
    private UserRepository userRepository;
    @Autowired
    private MediaRepository mediaRepository;
    @Autowired
    private RecommendationRepository recommendationRepository;
    private RecommendationService recommendationService;

    @BeforeEach
    public void initializeServices() {
        notificationService = new NotificationServiceMock();
        recommendationService = new RecommendationService(recommendationRepository, userRepository, mediaRepository, notificationService);
    }

    @Test
    public void checkThatRecommendationHasBeenSentWhenCreated() {
        var dummyDataCreator = new DummyDataCreator(userRepository, mediaRepository, recommendationRepository);
        var receiver = dummyDataCreator.generateUser(true);
        List<RecommendationDto> matchingRecommendations = new ArrayList<>();

        // create recommendation without existing user
        for (int i = 0; i < RandomUtils.nextInt(10, 30); i++) {
            var dummyMedia = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
            var dummySentBy = dummyDataCreator.generateUser(true);
            List<User> dummyReceivers = new ArrayList<>();
            for (int j = 0; j < RandomUtils.nextInt(3, 6); j++) {
                var dummyReceiver = dummyDataCreator.generateUser(true);
                dummyReceivers.add(dummyReceiver);
            }
            RecommendationCreateRequest creationRequest = new RecommendationCreateRequest(
                    dummyMedia.getId(),
                    dummyReceivers.stream().map(User::getId).toList(),
                    RandomStringUtils.random(20, true, false),
                    RandomUtils.nextInt(1, 5)
            );
            var result = recommendationService.createRecommendation(creationRequest, dummySentBy.getUsername());
            Assertions.assertNotNull(result);
        }

        // create recommendation with existing user
        for (int i = 0; i < RandomUtils.nextInt(3, 10); i++) {
            var dummyMedia = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
            var dummySentBy = dummyDataCreator.generateUser(true);
            List<User> dummyReceivers = new ArrayList<>();
            for (int j = 0; j < RandomUtils.nextInt(3, 6); j++) {
                var dummyReceiver = dummyDataCreator.generateUser(true);
                dummyReceivers.add(dummyReceiver);
            }
            dummyReceivers.add(receiver);
            RecommendationCreateRequest creationRequest = new RecommendationCreateRequest(
                    dummyMedia.getId(),
                    dummyReceivers.stream().map(User::getId).toList(),
                    RandomStringUtils.random(20, true, false),
                    RandomUtils.nextInt(1, 5)
            );
            var result = recommendationService.createRecommendation(creationRequest, dummySentBy.getUsername());
            matchingRecommendations.add(result);
            Assertions.assertNotNull(result);
        }

        // check that receiver has been received all recommendations
        var receivedRecommendations = notificationService.getRecommendationsSentForUser(receiver.getUsername());
        Assertions.assertEquals(receivedRecommendations.size(), matchingRecommendations.size());
        matchingRecommendations.forEach(expectedRecommendation -> {
            Assertions.assertTrue(receivedRecommendations.stream().anyMatch(expectedRecommendation::equals));
        });
    }
}
