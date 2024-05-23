package com.teasy.CineCircleApi.services.externals.mediaProviders;

import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.enums.MediaProviderEnum;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.externals.ExternalMedia;
import com.teasy.CineCircleApi.models.externals.ExternalMediaShort;
import org.springframework.data.domain.Pageable;

import java.util.List;

public interface MediaProvider {
    List<ExternalMediaShort> searchMedia(Pageable pageable, MediaSearchRequest mediaSearchRequest);

    ExternalMedia getMedia(String externalId, MediaTypeEnum mediaType) throws ExpectedException;

    List<String> listGenres();

    MediaProviderEnum getMediaProvider();

    List<String> getWatchProvidersForMedia(String externalId, MediaTypeEnum mediaType);
}
