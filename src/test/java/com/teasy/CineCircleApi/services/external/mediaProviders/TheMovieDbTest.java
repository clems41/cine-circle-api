package com.teasy.CineCircleApi.services.external.mediaProviders;


import com.teasy.CineCircleApi.config.RsaKeyProperties;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaType;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.services.NotificationService;
import com.teasy.CineCircleApi.services.RecommendationService;
import com.teasy.CineCircleApi.services.RecommendationServiceTest;
import com.teasy.CineCircleApi.services.externals.mediaProviders.theMovieDb.TheMovieDbService;
import com.teasy.CineCircleApi.services.utils.CustomHttpClient;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import org.apache.commons.lang3.RandomUtils;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.context.ActiveProfiles;

import java.util.*;

import static org.assertj.core.api.Assertions.assertThat;

@DataJpaTest
@ActiveProfiles("test")
public class TheMovieDbTest {
    @MockBean
    RsaKeyProperties rsaKeyProperties;
    @MockBean
    NotificationService notificationService;
    @MockBean
    CustomHttpClient httpClient;
    @Autowired
    private UserRepository userRepository;
    @Autowired
    private MediaRepository mediaRepository;
    @Autowired
    private RecommendationRepository recommendationRepository;
    private RecommendationService recommendationService;
    private TheMovieDbService theMovieDbService;

    @BeforeEach
    public void initializeServices() {
        recommendationService = new RecommendationService(recommendationRepository, userRepository, mediaRepository, notificationService);
        theMovieDbService = new TheMovieDbService(mediaRepository, httpClient, recommendationService);
    }

    @Test
    public void getMedia_CheckRecommendationFields() {
        // creation du user et du media en DB
        var dummyDataCreator = new DummyDataCreator(userRepository, mediaRepository, recommendationRepository);
        var user = dummyDataCreator.generateUser(true);
        var media = dummyDataCreator.generateMedia(true, MediaType.MOVIE);
        var wrongMedia = dummyDataCreator.generateMedia(true, MediaType.MOVIE);
        Set<User> receivers = new HashSet<>();
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            receivers.add(dummyDataCreator.generateUser(true));
        }
        receivers.add(user);
        Set<User> wrongReceivers = new HashSet<>();
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            receivers.add(dummyDataCreator.generateUser(true));
        }

        // creation de recommendations sur un autre media que celui créé en base où le user n'est pas en destinataire
        List<Recommendation> notMatchingRecommendations = new ArrayList<>();
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            notMatchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, null, wrongReceivers, wrongMedia));
        }

        // creation de recommendations sur un autre media que celui créé en base avec le user en destinataire
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            notMatchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, null, receivers, wrongMedia));
        }

        // creation de recommendations sur le media mais où le user n'est pas en destinataire
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            notMatchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, null, wrongReceivers, media));
        }

        // creation de recommendations sur le media avec le user en destinataire --> uniquement ces recommendations devront compter dans les champs du média
        List<Recommendation> matchingRecommendations = new ArrayList<>();
        for (int i = 0; i < RandomUtils.nextInt(5, 8); i++) {
            matchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, null, receivers, media));
        }

        // calcul des champs recommendations
        var expectedRecommendationCount = matchingRecommendations.size();
        var totalRating = matchingRecommendations
                .stream()
                .map(Recommendation::getRating)
                .reduce(0, Integer::sum);
        Double expectedRecommendationAverage = (double)totalRating / expectedRecommendationCount;

        // récupération du média et vérification des champs recommendations remplis
        var actualMedia = theMovieDbService.getMedia(media.getId(), user.getUsername());
        assertThat(actualMedia.getRecommendationRatingAverage()).isEqualTo(expectedRecommendationAverage);
        assertThat(actualMedia.getRecommendationRatingCount()).isEqualTo(expectedRecommendationCount);
        assertThat(actualMedia.getRecommendations().size()).isEqualTo(matchingRecommendations.size());
        // vérification que les recommendations correspondent à celle pour le média correspondant et le destinataire
        actualMedia.getRecommendations().forEach(mediaRecommendationDto -> {
            var expectedRecommendation = matchingRecommendations
                    .stream()
                    .filter(recommendation -> Objects.equals(recommendation.getId().toString(), mediaRecommendationDto.getId()))
                    .findAny();
            assertThat(expectedRecommendation.isPresent()).isTrue();
            assertThat(mediaRecommendationDto.getComment()).isEqualTo(expectedRecommendation.get().getComment());
            assertThat(mediaRecommendationDto.getRating()).isEqualTo(expectedRecommendation.get().getRating());
        });
    }
}
