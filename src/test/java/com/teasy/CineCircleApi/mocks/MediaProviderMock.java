package com.teasy.CineCircleApi.mocks;

import com.teasy.CineCircleApi.models.dtos.MediaCompleteDto;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import org.apache.commons.lang3.RandomStringUtils;
import org.springframework.data.domain.Pageable;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class MediaProviderMock implements MediaProvider {
    private final List<String> genres = new ArrayList<>();
    public MediaProviderMock() {
        for (int i = 0; i < 7; i++) {
            genres.add(RandomStringUtils.random(10, true, false));
        }
    }
    @Override
    public List<MediaDto> searchMedia(Pageable pageable, MediaSearchRequest mediaSearchRequest, String authenticatedUsername) {
        return null;
    }

    @Override
    public MediaCompleteDto getMedia(UUID internalId, String authenticatedUsername) {
        return null;
    }

    @Override
    public List<String> listGenres() {
        return genres;
    }
}
