package com.teasy.CineCircleApi.mocks;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.enums.MediaProviderEnum;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.models.externals.ExternalMedia;
import com.teasy.CineCircleApi.models.externals.ExternalMediaShort;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;
import org.springframework.data.domain.Pageable;

import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class MediaProviderMock implements MediaProvider {
    private final List<String> genres = new ArrayList<>();
    private final List<ExternalMedia> database;
    private final DummyDataCreator dummyDataCreator = new DummyDataCreator();

    public MediaProviderMock(List<ExternalMedia> database) {
        // create fake genres
        for (int i = 0; i < RandomUtils.nextInt(5, 10); i++) {
            genres.add(RandomStringUtils.random(10, true, false));
        }

        this.database = database;
    }

    @Override
    public List<ExternalMediaShort> searchMedia(Pageable pageable, MediaSearchRequest mediaSearchRequest) {
        List<ExternalMediaShort> result = new ArrayList<>();
        for (int i = 0; i < pageable.getPageSize(); i++) {
            result.add(fromExternalMediaToExternalMediaShort(generateMedia()));
        }
        return result;
    }

    @Override
    public ExternalMedia getMedia(String externalId, MediaTypeEnum mediaType) {
        if (database == null || database.isEmpty()) {
            var result = generateMedia();
            result.setExternalId(externalId);
            result.setMediaType(mediaType.name());
            return result;
        } else {
            return database
                    .stream()
                    .filter(externalMedia -> Objects.equals(externalMedia.getExternalId(), externalId))
                    .findAny()
                    .orElseThrow(() -> CustomExceptionHandler.mediaWithExternalIdNotFound(externalId));
        }
    }

    @Override
    public List<String> listGenres() {
        return genres;
    }

    @Override
    public MediaProviderEnum getMediaProvider() {
        return MediaProviderEnum.MOCK;
    }

    private ExternalMedia generateMedia() {
        var media = new ExternalMedia();
        media.setExternalId(String.valueOf(RandomUtils.nextInt(1_000, 100_000)));
        media.setMediaType(MediaTypeEnum.MOVIE.name());
        media.setTitle(RandomStringUtils.random(20, true, false));
        media.setOriginalTitle(RandomStringUtils.random(20, true, false));
        media.setPosterUrl(RandomStringUtils.random(20, true, true));
        media.setBackdropUrl(RandomStringUtils.random(20, true, true));
        media.setTrailerUrl(RandomStringUtils.random(20, true, true));
        media.setGenres(String.join(",",
                RandomStringUtils.random(6, true, false),
                RandomStringUtils.random(6, true, false)));
        media.setReleaseDate(LocalDate.now());
        media.setOverview(RandomStringUtils.random(100, true, false));
        media.setRuntime(RandomUtils.nextInt(40, 180));
        media.setOriginalLanguage(RandomStringUtils.random(6, true, false));
        media.setPopularity(RandomUtils.nextFloat(1, 10));
        media.setVoteAverage(RandomUtils.nextFloat(1, 10));
        media.setVoteCount(RandomUtils.nextInt(0, 1_000));
        media.setOriginCountry(RandomStringUtils.random(6, true, false));
        media.setActors(String.join(",",
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(10, true, false)));
        media.setDirector(RandomStringUtils.random(15, true, false));
        return media;
    }

    private ExternalMediaShort fromExternalMediaToExternalMediaShort(ExternalMedia media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, ExternalMediaShort.class);
    }
}
