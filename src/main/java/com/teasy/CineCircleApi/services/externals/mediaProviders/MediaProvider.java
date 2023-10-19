package com.teasy.CineCircleApi.services.externals.mediaProviders;

import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.dtos.requests.SearchMediaRequest;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

public interface MediaProvider {
    public Page<MediaShortDto> searchMedia(Pageable pageable, SearchMediaRequest searchMediaRequest);

    public Media getMedia(Long internalId);
}
