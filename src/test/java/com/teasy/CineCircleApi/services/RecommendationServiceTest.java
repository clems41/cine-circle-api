package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.*;
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
import java.util.UUID;

@DataJpaTest
@ActiveProfiles("test")
public class RecommendationServiceTest {
    NotificationServiceMock notificationService;
    @Autowired
    private UserRepository userRepository;
    @Autowired
    private MediaRepository mediaRepository;
    @Autowired
    private RecommendationRepository recommendationRepository;
    @Autowired
    private LibraryRepository libraryRepository;
    @Autowired
    private CircleRepository circleRepository;

    private LibraryService libraryService;
    private RecommendationService recommendationService;
    private DummyDataCreator dummyDataCreator;

    @BeforeEach
    public void setUp() {
        notificationService = new NotificationServiceMock();
        libraryService = new LibraryService(libraryRepository, mediaRepository, userRepository);
        recommendationService = new RecommendationService(recommendationRepository, userRepository, mediaRepository, libraryService, circleRepository, notificationService);
        dummyDataCreator = new DummyDataCreator(userRepository, mediaRepository, recommendationRepository, libraryRepository, circleRepository);
    }

    @Test
    public void checkThatRecommendationHasBeenSentWhenCreated() throws ExpectedException {
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
                    null,
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
                    null,
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

    @Test
    public void checkThatRecommendationHasBeenSentWhenCreated_PublicCircleUsers() throws ExpectedException {
        // In this testcase, we will send recommendation to a specific circle and check if all users from this circle have been received notification
        var circle = dummyDataCreator.generateCircle(true, null, true, null);

        // create recommendation without existing users from circle
        for (int i = 0; i < RandomUtils.nextInt(10, 30); i++) {
            var dummyMedia = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
            var dummySentBy = dummyDataCreator.generateUser(true);
            var dummyCircle = dummyDataCreator.generateCircle(true, null, true, null);
            List<UUID> circleIds = new ArrayList<>();
            circleIds.add(dummyCircle.getId());
            RecommendationCreateRequest creationRequest = new RecommendationCreateRequest(
                    dummyMedia.getId(),
                    null,
                    circleIds,
                    RandomStringUtils.random(20, true, false),
                    RandomUtils.nextInt(1, 5)
            );
            var result = recommendationService.createRecommendation(creationRequest, dummySentBy.getUsername());
            Assertions.assertNotNull(result);
        }

        // create one recommendation for existing circle
        var dummyMedia = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var dummySentBy = dummyDataCreator.generateUser(true);
        List<UUID> circleIds = new ArrayList<>();
        circleIds.add(circle.getId());
        RecommendationCreateRequest creationRequest = new RecommendationCreateRequest(
                dummyMedia.getId(),
                null,
                circleIds,
                RandomStringUtils.random(15, true, false),
                RandomUtils.nextInt(1, 5)
        );
        var actualRecommendation = recommendationService.createRecommendation(creationRequest, dummySentBy.getUsername());
        Assertions.assertNotNull(actualRecommendation);

        // check that all circle users have been received recommendation
        circle.getUsers().forEach(user -> {
            var receivedRecommendations = notificationService.getRecommendationsSentForUser(user.getUsername());
            Assertions.assertEquals(receivedRecommendations.size(), 1);
            Assertions.assertEquals(receivedRecommendations.get(0), actualRecommendation);
            // check that recommendation have been stored in database
            var recommendation = recommendationRepository.findById(UUID.fromString(receivedRecommendations.get(0).getId()));
            Assertions.assertTrue(recommendation.isPresent());
            Assertions.assertEquals(recommendation.get().getMedia().getId(), creationRequest.mediaId());
            Assertions.assertEquals(recommendation.get().getRating(), creationRequest.rating());
            Assertions.assertEquals(recommendation.get().getComment(), creationRequest.comment());
        });
    }
}
