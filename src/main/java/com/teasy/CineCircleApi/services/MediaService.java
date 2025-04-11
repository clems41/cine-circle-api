package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaFullDto;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.externals.ExternalMediaShort;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.utils.CustomExampleMatcher;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
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

    @Autowired
    public MediaService(MediaProvider theMovieDbService,
                        MediaRepository mediaRepository) {
        this.mediaProvider = theMovieDbService;
        this.mediaRepository = mediaRepository;
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
                log.error("Error while converting media entity to dto", e);
            }
        });
        return result;
    }

    public List<String> listGenres() {
        return mediaProvider.listGenres();
    }

    public List<String> getWatchProviders(UUID id) throws ExpectedException {
        var media = findMediaByIdOrElseThrow(id);
        return mediaProvider.getWatchProvidersForMedia(media.getExternalId(), media.getMediaType());
    }

    public MediaFullDto getMedia(UUID id, String authenticatedUsername) throws ExpectedException {
        // get media from database
        var media = findMediaByIdOrElseThrow(id);

        // complete it with more info if needed
        if (!media.getCompleted()) {
            completeMedia(media);
            mediaRepository.save(media);
        }
        return fromMediaEntityToDto(media, MediaFullDto.class, authenticatedUsername);
    }

    public Media findMediaByIdOrElseThrow(UUID id) throws ExpectedException {
        return mediaRepository
                .findById(id).
                orElseThrow(() -> new ExpectedException(ErrorDetails.ERR_MEDIA_NOT_FOUND.addingArgs(id)));
    }

    public <T> T fromMediaEntityToDto(Media media, Class<T> toValueType, String authenticatedUsername) throws ExpectedException {
        if (toValueType != MediaShortDto.class && toValueType != MediaFullDto.class) {
            throw new ExpectedException(ErrorDetails.ERR_MEDIA_TYPE_NOT_SUPPORTED.addingArgs(toValueType));
        }
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        var dto = mapper.convertValue(media, toValueType);
        if (toValueType == MediaShortDto.class) {
            MediaShortDto mediaDto = (MediaShortDto) dto;
            completeMediaShortDto(mediaDto, authenticatedUsername);
            return (T) mediaDto;
        } else {
            MediaFullDto mediaFullDto = (MediaFullDto) dto;
            completeMediaShortDto(mediaFullDto, authenticatedUsername);
            return (T) mediaFullDto;
        }
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

    private Optional<Media> findMediaWithExternalId(String externalId) {
        // build example matcher with external id and media provider
        ExampleMatcher matcher = CustomExampleMatcher.matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setExternalId(String.valueOf(externalId));
        exampleMedia.setMediaProvider(mediaProvider.getMediaProvider().name());

        return mediaRepository
                .findOne(Example.of(exampleMedia, matcher));
    }

    private void completeMediaShortDto(MediaShortDto mediaShortDto, String authenticatedUsername) {
        mediaShortDto.setIsInLibrary(mediaRepository.isInLibraryOfUser(mediaShortDto.getId(), authenticatedUsername));
        mediaShortDto.setIsInWatchlist(mediaRepository.isInWatchlistOfUser(mediaShortDto.getId(), authenticatedUsername));
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
