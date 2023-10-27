package com.teasy.CineCircleApi.utils;

import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaProvider;
import com.teasy.CineCircleApi.models.enums.MediaType;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.HashSet;
import java.util.Set;

@NoArgsConstructor
@AllArgsConstructor
public class DummyDataCreator {
    private UserRepository userRepository;
    private MediaRepository mediaRepository;
    private RecommendationRepository recommendationRepository;

    public User generateUser(Boolean storeInDatabase) {
        var displayName = RandomStringUtils.random(20, true, false);
        var username = RandomStringUtils.random(10, true, true);
        var email = String.format("%s@%s.com",
                RandomStringUtils.random(10, true, true),
                RandomStringUtils.random(6, true, false));
        var hashPassword = RandomStringUtils.random(16, true, true);
        var user = new User(username, email, hashPassword, displayName);
        if (userRepository != null && storeInDatabase) {
            return userRepository.save(user);
        }
        return user;
    }

    public Media generateMedia(Boolean storeInDatabase, MediaType mediaType) {
        var media = new Media();
        media.setExternalId(String.valueOf(RandomUtils.nextInt(1_000, 100_000)));
        media.setMediaProvider(MediaProvider.THE_MOVIE_DATABASE.name());
        media.setTitle(RandomStringUtils.random(20, true, false));
        media.setOriginalTitle(RandomStringUtils.random(20, true, false));
        media.setPosterUrl(RandomStringUtils.random(20, true, true));
        media.setBackdropUrl(RandomStringUtils.random(20, true, true));
        media.setTrailerUrl(RandomStringUtils.random(20, true, true));
        media.setGenres(String.join(",",
                RandomStringUtils.random(6, true, false),
                RandomStringUtils.random(6, true, false)));
        if (mediaType != null) {
            media.setMediaType(mediaType.name());
        } else {
            media.setMediaType(MediaType.MOVIE.name());
        }
        media.setReleaseDate(LocalDate.now());
        media.setRuntime(RandomUtils.nextInt(40, 180));
        media.setOriginalLanguage(RandomStringUtils.random(6, true, false));
        media.setPopularity(RandomUtils.nextFloat(1, 10));
        media.setVoteAverage(RandomUtils.nextFloat(1, 10));
        media.setVoteCount(RandomUtils.nextInt(0, 1_000));
        media.setOriginCountry(RandomStringUtils.random(6, true, false));
        media.setCompleted(true);
        media.setActors(String.join(",",
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(10, true, false)));
        media.setDirector(RandomStringUtils.random(15, true, false));
        if (mediaRepository != null && storeInDatabase) {
            return mediaRepository.save(media);
        }
        return media;
    }

    public Recommendation generateRecommendation(Boolean storeInDatabase, User sentBy, Set<User> receivers, Media media) {
        var recommendation = new Recommendation();
        recommendation.setSentAt(LocalDateTime.now());
        recommendation.setComment(RandomStringUtils.random(30, true, false));
        recommendation.setRating(RandomUtils.nextInt(1, 5));
        if (sentBy != null) {
            recommendation.setSentBy(sentBy);
        } else {
            recommendation.setSentBy(generateUser(storeInDatabase));
        }
        if (receivers != null) {
            recommendation.setReceivers(receivers);
        } else {
            var receiversSize = RandomUtils.nextInt(1, 5);
            Set<User> generatedReceivers = new HashSet<>();
            for (int i = 0; i < receiversSize; i++) {
                generatedReceivers.add(generateUser(storeInDatabase));
            }
            recommendation.setReceivers(generatedReceivers);
        }
        if (media != null) {
            recommendation.setMedia(media);
        } else {
            recommendation.setMedia(generateMedia(storeInDatabase, null));
        }
        if (recommendationRepository != null && storeInDatabase) {
            return recommendationRepository.save(recommendation);
        }
        return recommendation;
    }

}
