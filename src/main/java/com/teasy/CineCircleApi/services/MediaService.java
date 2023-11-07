package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaCompleteDto;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.MediaRecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.models.externals.ExternalMediaShort;
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

    @Autowired
    public MediaService(MediaProvider theMovieDbService,
                        MediaRepository mediaRepository,
                        RecommendationRepository recommendationRepository,
                        UserRepository userRepository) {
        this.mediaProvider = theMovieDbService;
        this.mediaRepository = mediaRepository;
        this.recommendationRepository = recommendationRepository;
        this.userRepository = userRepository;
    }

    public List<MediaDto> searchMedia(Pageable pageable, MediaSearchRequest mediaSearchRequest, String authenticatedUsername) {
        var medias = mediaProvider.searchMedia(pageable, mediaSearchRequest);
        List<MediaDto> result = new ArrayList<>();
        medias.forEach(externalMedia -> {
            // store result in database with internalId if not already exists
            var existingMedia = findMediaWithExternalId(externalMedia.getExternalId());
            if (existingMedia.isEmpty()) {
                var newMedia = fromExternalMediaShortToMediaEntity(externalMedia);
                newMedia.setCompleted(false);
                mediaRepository.save(newMedia);
                result.add(fromMediaEntityToDto(newMedia, MediaDto.class, authenticatedUsername));
            } else {
                result.add(fromMediaEntityToDto(existingMedia.get(), MediaDto.class, authenticatedUsername));
            }
        });
        return result;
    }

    public List<String> listGenres() {
        return mediaProvider.listGenres();
    }

    public MediaCompleteDto getMedia(UUID id, String authenticatedUsername) throws CustomException {
        // build example matcher with id
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setId(id);

        // get media from database
        var media = mediaRepository
                .findOne(Example.of(exampleMedia, matcher))
                .orElseThrow(() -> CustomExceptionHandler.mediaWithIdNotFound(id));

        // complete it with more info if needed
        if (!media.getCompleted()) {
            completeMedia(media);
            mediaRepository.save(media);
        }
        return fromMediaEntityToDto(media, MediaCompleteDto.class, authenticatedUsername);
    }

    private void completeMedia(Media media) {
        var completedMedia = mediaProvider.getMedia(media.getExternalId(), MediaTypeEnum.valueOf(media.getMediaType()));
        media.setDirector(completedMedia.getDirector());
        media.setActors(completedMedia.getActors());
        media.setTrailerUrl(completedMedia.getTrailerUrl());
        media.setGenres(completedMedia.getGenres());
        media.setVoteAverage(completedMedia.getVoteAverage());
        media.setVoteCount(completedMedia.getVoteCount());
        media.setPopularity(completedMedia.getPopularity());
        media.setOriginCountry(completedMedia.getOriginCountry());
        media.setCompleted(true);
    }

    private <T> void addRecommendationRatingFields(T mediaDto, String authenticatedUsername) {
        if (mediaDto.getClass() != MediaDto.class && mediaDto.getClass() != MediaCompleteDto.class) {
            return;
        }

        // find recommendation average and count
        UUID mediaId;
        if (mediaDto.getClass() == MediaDto.class) {
            mediaId = UUID.fromString(((MediaDto) mediaDto).getId());
        } else {
            mediaId = UUID.fromString(((MediaCompleteDto) mediaDto).getId());
        }
        var recommendations = findRecommendationsReceivedForMediaAndAuthenticatedUsername(mediaId, authenticatedUsername);
        var recommendationRatingCount = recommendations.size();
        var recommendationRatingAverage = recommendations
                .stream()
                .mapToDouble(Recommendation::getRating)
                .average();
        if (mediaDto.getClass() == MediaDto.class) {
            ((MediaDto) mediaDto).setRecommendationRatingCount(recommendationRatingCount);
            ((MediaDto) mediaDto).setRecommendationRatingAverage(recommendationRatingAverage.isPresent() ?
                    recommendationRatingAverage.getAsDouble() : null);
        } else {
            ((MediaCompleteDto) mediaDto).setRecommendationRatingCount(recommendationRatingCount);
            ((MediaCompleteDto) mediaDto).setRecommendationRatingAverage(recommendationRatingAverage.isPresent() ?
                    recommendationRatingAverage.getAsDouble() : null);
            // add all recommendations comment for complete media dto
            ((MediaCompleteDto) mediaDto).setRecommendations(
                    recommendations.stream().map(this::fromRecommendationToMediaRecommendationDto).toList()
            );
        }
    }

    private List<Recommendation> findRecommendationsReceivedForMediaAndAuthenticatedUsername(UUID mediaId, String username) {
        var user = findUserByUsernameOrElseThrow(username);
        return recommendationRepository.findAllByReceivers_IdAndMedia_Id(
                        PageRequest.ofSize(1000),
                        user.getId(),
                        mediaId
                )
                .stream()
                .toList();
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

    private User findUserByUsernameOrElseThrow(String username) throws CustomException {
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
    }

    private <T> T fromMediaEntityToDto(Media media, Class<T> toValueType, String authenticatedUsername) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        var result = mapper.convertValue(media, toValueType);
        addRecommendationRatingFields(result, authenticatedUsername);
        return result;
    }

    private MediaRecommendationDto fromRecommendationToMediaRecommendationDto(Recommendation recommendationDto) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(recommendationDto, MediaRecommendationDto.class);
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
