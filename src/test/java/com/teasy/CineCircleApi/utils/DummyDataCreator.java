package com.teasy.CineCircleApi.utils;

import com.teasy.CineCircleApi.models.entities.*;
import com.teasy.CineCircleApi.models.enums.MediaProviderEnum;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.repositories.*;
import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.Objects;

@NoArgsConstructor
@AllArgsConstructor
public class DummyDataCreator {
    private UserRepository userRepository;
    private MediaRepository mediaRepository;
    private RecommendationRepository recommendationRepository;
    private LibraryRepository libraryRepository;
    private CircleRepository circleRepository;

    public User generateUser(Boolean storeInDatabase) {
        var displayName = RandomUtils.randomString(20);
        var username = RandomUtils.randomString(10);
        var email = String.format("%s@%s.com",
                RandomUtils.randomString(10),
                RandomUtils.randomString(6));
        var hashPassword = RandomUtils.randomString(16);
        var user = new User(username, email, hashPassword, displayName);
        if (userRepository != null && storeInDatabase) {
            return userRepository.save(user);
        }
        return user;
    }

    public Media generateMedia(Boolean storeInDatabase, MediaTypeEnum mediaTypeEnum) {
        var media = new Media();
        media.setExternalId(String.valueOf(RandomUtils.randomInt(1_000, 999_000)));
        media.setMediaProvider(MediaProviderEnum.THE_MOVIE_DATABASE.name());
        media.setTitle(RandomUtils.randomString(20));
        media.setOriginalTitle(RandomUtils.randomString(20));
        media.setPosterUrl(RandomUtils.randomString(20));
        media.setBackdropUrl(RandomUtils.randomString(20));
        media.setTrailerUrl(RandomUtils.randomString(20));
        media.setGenres(String.join(",",
                RandomUtils.randomString(6),
                RandomUtils.randomString(6)));
        media.setMediaType(Objects.requireNonNullElse(mediaTypeEnum, MediaTypeEnum.MOVIE));
        media.setOverview(RandomUtils.randomString(100));
        media.setReleaseDate(LocalDate.now());
        media.setRuntime(RandomUtils.randomInt(40, 180));
        media.setOriginalLanguage(RandomUtils.randomString(6));
        media.setPopularity(RandomUtils.randomFloat(1, 10));
        media.setVoteAverage(RandomUtils.randomFloat(1, 10));
        media.setVoteCount(RandomUtils.randomInt(0, 1_000));
        media.setOriginCountry(RandomUtils.randomString(6));
        media.setCompleted(true);
        media.setActors(String.join(",",
                RandomUtils.randomString(10),
                RandomUtils.randomString(10),
                RandomUtils.randomString(10),
                RandomUtils.randomString(10)));
        media.setDirector(RandomUtils.randomString(15));
        if (mediaRepository != null && storeInDatabase) {
            return mediaRepository.save(media);
        }
        return media;
    }

    public void generateRecommendation(Boolean storeInDatabase, User sentBy, User receiver, Media media, Boolean read) {
        var recommendation = new Recommendation();
        recommendation.setSentAt(LocalDateTime.now());
        recommendation.setComment(RandomUtils.randomString(30));
        recommendation.setRating(RandomUtils.randomInt(1, 5));
        if (sentBy != null) {
            recommendation.setSentBy(sentBy);
        } else {
            recommendation.setSentBy(generateUser(storeInDatabase));
        }
        if (receiver != null) {
            recommendation.setReceiver(receiver);
        } else {
            recommendation.setReceiver(generateUser(storeInDatabase));
        }
        if (media != null) {
            recommendation.setMedia(media);
        } else {
            recommendation.setMedia(generateMedia(storeInDatabase, null));
        }
        if (read != null) {
            recommendation.setRead(read);
        }
        if (recommendationRepository != null && storeInDatabase) {
            recommendationRepository.save(recommendation);
        }
    }

    public void addMediaToLibrary(User user, Media media) {
        var libraryRecord = new Library(
                user,
                media,
                RandomUtils.randomString(20),
                RandomUtils.randomInt(1, 5));
        libraryRepository.save(libraryRecord);
    }

    public Circle generateCircle(Boolean storeInDatabase, User creator, Boolean isPublic, String name) {
        var description = RandomUtils.randomString(30);
        if (name == null) {
            name = RandomUtils.randomString(15);
        }
        if (creator == null) {
            creator = generateUser(storeInDatabase);
        }
        if (isPublic == null) {
            isPublic = false;
        }
        var circle = new Circle(isPublic, name, description, creator);
        for (int i = 0; i < RandomUtils.randomInt(2, 8); i++) {
            circle.addUser(generateUser(storeInDatabase));
        }
        if(storeInDatabase) {
            circleRepository.save(circle);
        }
        return circle;
    }

}
