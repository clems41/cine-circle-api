package com.teasy.CineCircleApi.services.externals.mediaProviders;

import com.teasy.CineCircleApi.models.dtos.MediaCompleteDto;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.dtos.responses.MediaGenreResponse;
import org.springframework.data.domain.Pageable;

import java.util.List;

public interface MediaProvider {
    List<MediaDto> searchMedia(Pageable pageable,
                               MediaSearchRequest mediaSearchRequest,
                               String authenticatedUsername);

    MediaCompleteDto getMedia(Long internalId,
                              String authenticatedUsername);

    MediaGenreResponse listGenres();
}
