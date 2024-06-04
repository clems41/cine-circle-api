package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaFullDto;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.RecommendationMediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.entities.Library;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.ErrorMessage;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.externals.ExternalMediaShort;
import com.teasy.CineCircleApi.repositories.LibraryRepository;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
@Slf4j
public class MediaService {
    private final MediaProvider mediaProvider;

    private final MediaRepository mediaRepository;
    private final RecommendationRepository recommendationRepository;
    private final UserRepository userRepository;
    private final LibraryRepository libraryRepository;

    @Autowired
    public MediaService(MediaProvider theMovieDbService,
                        MediaRepository mediaRepository,
                        RecommendationRepository recommendationRepository,
                        UserRepository userRepository,
                        LibraryRepository libraryRepository) {
        this.mediaProvider = theMovieDbService;
        this.mediaRepository = mediaRepository;
        this.recommendationRepository = recommendationRepository;
        this.userRepository = userRepository;
        this.libraryRepository = libraryRepository;
    }

    public List<MediaShortDto> searchMedia(Pageable pageable, MediaSearchRequest mediaSearchRequest, String authenticatedUsername) {
        var medias = mediaProvider.searchMedia(pageable, mediaSearchRequest);
        List<MediaShortDto> result = new ArrayList<>();
        medias.forEach(externalMedia -> {
            // store result in database with internalId if not already exists
            var existingMedia = findMediaWithExternalId(externalMedia.getExternalId());
            try {
                if (existingMedia.isPresent()) {
                    result.add(fromMediaEntityToDto(existingMedia.get(), MediaShortDto.class, authenticatedUsername));
                } else {
                    var newMedia = fromExternalMediaShortToMediaEntity(externalMedia);
                    newMedia.setCompleted(false);
                    mediaRepository.save(newMedia);
                    result.add(fromMediaEntityToDto(newMedia, MediaShortDto.class, authenticatedUsername));
                }
            } catch (ExpectedException e) {
                log.error("Error while converting media entity to dto: " + e.getMessage());
            }
        });
        return result;
    }

    public List<String> listGenres() {
        return mediaProvider.listGenres();
    }

    public List<String> getWatchProviders(UUID id) throws ExpectedException {
        var media = findMediaWithIdOrElseThrow(id);
        return mediaProvider.getWatchProvidersForMedia(media.getExternalId(), media.getMediaType());
    }

    public MediaFullDto getMedia(UUID id, String authenticatedUsername) throws ExpectedException {
        // get media from database
        var media = findMediaWithIdOrElseThrow(id);

        // complete it with more info if needed
        if (!media.getCompleted()) {
            completeMedia(media);
            mediaRepository.save(media);
        }
        return fromMediaEntityToDto(media, MediaFullDto.class, authenticatedUsername);
    }

    private void completeMedia(Media media) throws ExpectedException {
        var completedMedia = mediaProvider.getMedia(media.getExternalId(), media.getMediaType());
        media.setDirector(completedMedia.getDirector());
        media.setActors(completedMedia.getActors());
        media.setTrailerUrl(completedMedia.getTrailerUrl());
        media.setGenres(completedMedia.getGenres());
        media.setVoteAverage(completedMedia.getVoteAverage());
        media.setVoteCount(completedMedia.getVoteCount());
        media.setPopularity(completedMedia.getPopularity());
        media.setOriginCountry(completedMedia.getOriginCountry());
        media.setRuntime(completedMedia.getRuntime());
        media.setCompleted(true);
    }

    private <T> void addRecommendationRatingFields(T mediaDto, String authenticatedUsername) throws ExpectedException {
        if (mediaDto.getClass() != MediaShortDto.class && mediaDto.getClass() != MediaFullDto.class) {
            return;
        }

        // find recommendation average and count
        UUID mediaId;
        if (mediaDto.getClass() == MediaShortDto.class) {
            mediaId = UUID.fromString(((MediaShortDto) mediaDto).getId());
        } else {
            mediaId = UUID.fromString(((MediaFullDto) mediaDto).getId());
        }
        var recommendationsReceived = findRecommendationsReceivedForMediaAndAuthenticatedUsername(mediaId, authenticatedUsername);
        var recommendationRatingCount = recommendationsReceived.size();
        var recommendationRatingAverage = recommendationsReceived
                .stream()
                .mapToDouble(Recommendation::getRating)
                .average();
        if (mediaDto.getClass() == MediaShortDto.class) {
            ((MediaShortDto) mediaDto).setRecommendationRatingCount(recommendationRatingCount);
            ((MediaShortDto) mediaDto).setRecommendationRatingAverage(recommendationRatingAverage.isPresent() ?
                    recommendationRatingAverage.getAsDouble() : null);
        } else {
            ((MediaFullDto) mediaDto).setRecommendationRatingCount(recommendationRatingCount);
            ((MediaFullDto) mediaDto).setRecommendationRatingAverage(recommendationRatingAverage.isPresent() ?
                    recommendationRatingAverage.getAsDouble() : null);
            // add all recommendations received for complete media dto
            ((MediaFullDto) mediaDto).setRecommendationsReceived(
                    recommendationsReceived.stream().map(this::fromRecommendationToMediaRecommendationDto).toList()
            );
            // add all recommendations sent for complete media dto
            var recommendationsSent = findRecommendationsSentForMediaAndAuthenticatedUsername(mediaId, authenticatedUsername);
            ((MediaFullDto) mediaDto).setRecommendationsSent(
                    recommendationsSent.stream().map(this::fromRecommendationToMediaRecommendationDto).toList()
            );
        }
    }

    private <T> void addPersonalFields(T mediaDto, String authenticatedUsername) throws ExpectedException {
        if (mediaDto.getClass() != MediaShortDto.class && mediaDto.getClass() != MediaFullDto.class) {
            return;
        }

        // extract media Id from dto
        UUID mediaId;
        if (mediaDto.getClass() == MediaShortDto.class) {
            mediaId = UUID.fromString(((MediaShortDto) mediaDto).getId());
        } else {
            mediaId = UUID.fromString(((MediaFullDto) mediaDto).getId());
        }

        // find user and media
        var user = findUserByUsernameOrElseThrow(authenticatedUsername);
        var media = mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.MEDIA_NOT_FOUND, HttpStatus.NOT_FOUND));

        // find library record if exist
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingLibrary = new Library(user, media, null, null);
        matchingLibrary.setAddedAt(null);
        var libraryRecord = libraryRepository.findOne(Example.of(matchingLibrary, matcher));
        if (libraryRecord.isPresent()) {
            if (mediaDto.getClass() == MediaShortDto.class) {
                ((MediaShortDto) mediaDto).setPersonalRating(libraryRecord.get().getRating());
                ((MediaShortDto) mediaDto).setPersonalComment(libraryRecord.get().getComment());
            } else {
                ((MediaFullDto) mediaDto).setPersonalRating(libraryRecord.get().getRating());
                ((MediaFullDto) mediaDto).setPersonalComment(libraryRecord.get().getComment());
            }
        }
    }

    private List<Recommendation> findRecommendationsReceivedForMediaAndAuthenticatedUsername(UUID mediaId, String username) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        return recommendationRepository.findAllByReceivers_IdAndMedia_Id(
                        PageRequest.ofSize(1000),
                        user.getId(),
                        mediaId
                )
                .stream()
                .toList();
    }

    private List<Recommendation> findRecommendationsSentForMediaAndAuthenticatedUsername(UUID mediaId, String username) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        var matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        User matchingSender = new User();
        matchingSender.setId(user.getId());
        Media matchingMedia = new Media();
        matchingMedia.setId(mediaId);
        var matchingRecommendation = new Recommendation(matchingSender, matchingMedia, null, null, null, null);
        matchingRecommendation.setSentAt(null);

        return recommendationRepository.findAll(Example.of(matchingRecommendation, matcher));
    }

    private Optional<Media> findMediaWithExternalId(String externalId) {
        // build example matcher with external id and media provider
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setExternalId(String.valueOf(externalId));
        exampleMedia.setMediaProvider(mediaProvider.getMediaProvider().name());

        return mediaRepository
                .findOne(Example.of(exampleMedia, matcher));
    }

    private Media findMediaWithIdOrElseThrow(UUID id) throws ExpectedException {
        return mediaRepository
                .findById(id).
                orElseThrow(() -> new ExpectedException(ErrorMessage.MEDIA_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private User findUserByUsernameOrElseThrow(String username) throws ExpectedException {
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ExpectedException(ErrorMessage.USER_NOT_FOUND, HttpStatus.NOT_FOUND));
    }

    private <T> T fromMediaEntityToDto(Media media, Class<T> toValueType, String authenticatedUsername) throws ExpectedException {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        var result = mapper.convertValue(media, toValueType);
        addRecommendationRatingFields(result, authenticatedUsername);
        addPersonalFields(result, authenticatedUsername);
        return result;
    }

    private RecommendationMediaDto fromRecommendationToMediaRecommendationDto(Recommendation recommendationDto) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(recommendationDto, RecommendationMediaDto.class);
    }

    private Media fromExternalMediaShortToMediaEntity(ExternalMediaShort externalMediaShort) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        var result = mapper.convertValue(externalMediaShort, Media.class);
        result.setMediaProvider(mediaProvider.getMediaProvider().name());
        return result;
    }
}
