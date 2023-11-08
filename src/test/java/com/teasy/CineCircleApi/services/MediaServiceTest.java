package com.teasy.CineCircleApi.services;


import com.teasy.CineCircleApi.config.RsaKeyProperties;
import com.teasy.CineCircleApi.mocks.MediaProviderMock;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.repositories.*;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.utils.CustomHttpClient;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import org.apache.commons.lang3.RandomUtils;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.data.domain.PageRequest;
import org.springframework.test.context.ActiveProfiles;

import java.util.*;

import static org.assertj.core.api.Assertions.assertThat;

@DataJpaTest
@ActiveProfiles("test")
public class MediaServiceTest {
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
    @Autowired
    private LibraryRepository libraryRepository;
    @Autowired
    private CircleRepository circleRepository;
    private MediaService mediaService;
    private DummyDataCreator dummyDataCreator;

    @BeforeEach
    public void initializeServices() {
        MediaProvider mediaProvider = new MediaProviderMock(new ArrayList<>());
        mediaService = new MediaService(mediaProvider, mediaRepository, recommendationRepository, userRepository, libraryRepository);
        dummyDataCreator = new DummyDataCreator(userRepository, mediaRepository, recommendationRepository, libraryRepository, circleRepository);
    }

    @Test
    public void getMedia_CheckRecommendationFields_RecommendationsReceived() {
        // creation du user et du media en DB
        var user = dummyDataCreator.generateUser(true);
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var wrongMedia = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
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
        var actualMedia = mediaService.getMedia(media.getId(), user.getUsername());
        assertThat(actualMedia.getRecommendationRatingAverage()).isEqualTo(expectedRecommendationAverage);
        assertThat(actualMedia.getRecommendationRatingCount()).isEqualTo(expectedRecommendationCount);
        assertThat(actualMedia.getRecommendationsReceived().size()).isEqualTo(matchingRecommendations.size());
        // vérification que les recommendations reçues correspondent à celles pour le média correspondant et le destinataire
        actualMedia.getRecommendationsReceived().forEach(recommendationMediaDto -> {
            var expectedRecommendation = matchingRecommendations
                    .stream()
                    .filter(recommendation -> Objects.equals(recommendation.getId().toString(), recommendationMediaDto.getId()))
                    .findAny();
            assertThat(expectedRecommendation.isPresent()).isTrue();
            assertThat(recommendationMediaDto.getComment()).isEqualTo(expectedRecommendation.get().getComment());
            assertThat(recommendationMediaDto.getRating()).isEqualTo(expectedRecommendation.get().getRating());
        });
    }

    @Test
    public void getMedia_CheckRecommendationFields_RecommendationsSent() {
        // creation du user et du media en DB
        var sender = dummyDataCreator.generateUser(true);
        var wrongSender = dummyDataCreator.generateUser(true);
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var wrongMedia = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);

        // creation de recommendations sur un autre media que celui créé en base où le user n'est pas l'envoyeur
        List<Recommendation> notMatchingRecommendations = new ArrayList<>();
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            notMatchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, wrongSender, null, wrongMedia));
        }

        // creation de recommendations sur un autre media que celui créé en base avec le user est l'envoyeur
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            notMatchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, sender, null, wrongMedia));
        }

        // creation de recommendations sur le media mais où le user n'est pas l'envoyeur
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) {
            notMatchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, wrongSender, null, media));
        }

        // creation de recommendations sur le media avec le user est l'envoyeur --> uniquement ces recommendations devront compter dans les champs du média
        List<Recommendation> matchingRecommendations = new ArrayList<>();
        for (int i = 0; i < RandomUtils.nextInt(5, 8); i++) {
            matchingRecommendations.add(dummyDataCreator.generateRecommendation(
                    true, sender, null, media));
        }

        // récupération du média et vérification des champs recommendations remplis
        var actualMedia = mediaService.getMedia(media.getId(), sender.getUsername());
        assertThat(actualMedia.getRecommendationsSent().size()).isEqualTo(matchingRecommendations.size());
        // vérification que les recommendations reçues correspondent à celles pour le média correspondant et le destinataire
        actualMedia.getRecommendationsReceived().forEach(recommendationMediaDto -> {
            var expectedRecommendation = matchingRecommendations
                    .stream()
                    .filter(recommendation -> Objects.equals(recommendation.getId().toString(), recommendationMediaDto.getId()))
                    .findAny();
            assertThat(expectedRecommendation.isPresent()).isTrue();
            assertThat(recommendationMediaDto.getComment()).isEqualTo(expectedRecommendation.get().getComment());
            assertThat(recommendationMediaDto.getRating()).isEqualTo(expectedRecommendation.get().getRating());
        });
    }

    @Test
    public void searchThenGetMedia_ShouldBeCompleted() {
        /* Create user */
        var user = dummyDataCreator.generateUser(true);

        /* Search some medias */
        var medias = mediaService.searchMedia(
                PageRequest.ofSize(10),
                new MediaSearchRequest("inception"),
                user.getUsername());

        /* Check that medias have been created in database and not complete */
        medias.forEach(mediaDto -> {
            var media = mediaRepository.findById(UUID.fromString(mediaDto.getId()));
            assertThat(media.isPresent()).isTrue();

            // check media entity fields
            assertThat(media.get().getId()).isNotNull();
            assertThat(media.get().getTitle()).isNotEmpty();
            assertThat(media.get().getOriginalTitle()).isNotEmpty();
            assertThat(media.get().getPosterUrl()).isNotEmpty();
            assertThat(media.get().getBackdropUrl()).isNotEmpty();
            assertThat(media.get().getOverview()).isNotEmpty();
            assertThat(media.get().getMediaType()).isNotEmpty();
            assertThat(media.get().getReleaseDate()).isNotNull();
            assertThat(media.get().getRuntime()).isNotZero().isPositive();
            assertThat(media.get().getOriginalLanguage()).isNotEmpty();
            assertThat(media.get().getCompleted()).isFalse();

            // check mediaDto fields
            assertThat(mediaDto.getId()).isNotEmpty();
            assertThat(mediaDto.getTitle()).isNotEmpty();
            assertThat(mediaDto.getOriginalTitle()).isNotEmpty();
            assertThat(mediaDto.getPosterUrl()).isNotEmpty();
            assertThat(mediaDto.getBackdropUrl()).isNotEmpty();
            assertThat(mediaDto.getOverview()).isNotEmpty();
            assertThat(mediaDto.getMediaType()).isNotEmpty();
            assertThat(mediaDto.getReleaseDate()).isNotNull();
            assertThat(mediaDto.getRuntime()).isNotZero().isPositive();
            assertThat(mediaDto.getOriginalLanguage()).isNotEmpty();
            assertThat(mediaDto.getRecommendationRatingCount()).isZero();
            assertThat(mediaDto.getRecommendationRatingAverage()).isNull();
        });

        /* Get one of these media and check if all fields have been filled (including new ones : genres, actors, director, trailer, etc.) */
        var expectedMedia = medias.get(4);
        var actualMedia = mediaService.getMedia(UUID.fromString(expectedMedia.getId()), user.getUsername());
        assertThat(actualMedia.getId()).isEqualTo(expectedMedia.getId());
        assertThat(actualMedia.getTitle()).isEqualTo(expectedMedia.getTitle());
        assertThat(actualMedia.getOriginalTitle()).isEqualTo(expectedMedia.getOriginalTitle());
        assertThat(actualMedia.getPosterUrl()).isEqualTo(expectedMedia.getPosterUrl());
        assertThat(actualMedia.getBackdropUrl()).isEqualTo(expectedMedia.getBackdropUrl());
        assertThat(actualMedia.getOverview()).isEqualTo(expectedMedia.getOverview());
        assertThat(actualMedia.getMediaType()).isEqualTo(expectedMedia.getMediaType());
        assertThat(actualMedia.getReleaseDate()).isEqualTo(expectedMedia.getReleaseDate());
        assertThat(actualMedia.getRuntime()).isEqualTo(expectedMedia.getRuntime());
        assertThat(actualMedia.getOriginalLanguage()).isEqualTo(expectedMedia.getOriginalLanguage());
        assertThat(actualMedia.getRecommendationRatingCount()).isEqualTo(expectedMedia.getRecommendationRatingCount());
        assertThat(actualMedia.getRecommendationRatingAverage()).isEqualTo(expectedMedia.getRecommendationRatingAverage());
        assertThat(actualMedia.getGenres()).isNotEmpty();
        assertThat(actualMedia.getActors()).isNotEmpty();
        assertThat(actualMedia.getDirector()).isNotEmpty();
        assertThat(actualMedia.getTrailerUrl()).isNotEmpty();
        assertThat(actualMedia.getPopularity()).isNotZero();
        assertThat(actualMedia.getVoteAverage()).isNotZero();
        assertThat(actualMedia.getVoteCount()).isNotZero();
        assertThat(actualMedia.getOriginCountry()).isNotEmpty();

        /* Check that media in database is complete with all fields filled */
        var actualMediaEntity = mediaRepository.findById(UUID.fromString(actualMedia.getId()));
        assertThat(actualMediaEntity.isPresent()).isTrue();
        assertThat(actualMediaEntity.get().getId().toString()).isEqualTo(actualMedia.getId());
        assertThat(actualMediaEntity.get().getTitle()).isEqualTo(actualMedia.getTitle());
        assertThat(actualMediaEntity.get().getOriginalTitle()).isEqualTo(actualMedia.getOriginalTitle());
        assertThat(actualMediaEntity.get().getPosterUrl()).isEqualTo(actualMedia.getPosterUrl());
        assertThat(actualMediaEntity.get().getBackdropUrl()).isEqualTo(actualMedia.getBackdropUrl());
        assertThat(actualMediaEntity.get().getOverview()).isEqualTo(actualMedia.getOverview());
        assertThat(actualMediaEntity.get().getMediaType()).isEqualTo(actualMedia.getMediaType());
        assertThat(actualMediaEntity.get().getReleaseDate()).isEqualTo(actualMedia.getReleaseDate());
        assertThat(actualMediaEntity.get().getRuntime()).isEqualTo(actualMedia.getRuntime());
        assertThat(actualMediaEntity.get().getOriginalLanguage()).isEqualTo(actualMedia.getOriginalLanguage());
        assertThat(actualMediaEntity.get().getGenres()).isEqualTo(actualMedia.getGenres());
        assertThat(actualMediaEntity.get().getActors()).isEqualTo(actualMedia.getActors());
        assertThat(actualMediaEntity.get().getDirector()).isEqualTo(actualMedia.getDirector());
        assertThat(actualMediaEntity.get().getTrailerUrl()).isEqualTo(actualMedia.getTrailerUrl());
        assertThat(actualMediaEntity.get().getPopularity()).isEqualTo(actualMedia.getPopularity());
        assertThat(actualMediaEntity.get().getVoteAverage()).isEqualTo(actualMedia.getVoteAverage());
        assertThat(actualMediaEntity.get().getVoteCount()).isEqualTo(actualMedia.getVoteCount());
        assertThat(actualMediaEntity.get().getOriginCountry()).isEqualTo(actualMedia.getOriginCountry());
        assertThat(actualMediaEntity.get().getCompleted()).isTrue();
    }
}
