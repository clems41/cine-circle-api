package com.teasy.CineCircleApi.services.externals.mediaProviders;

import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.dtos.requests.SearchMediaRequest;
import org.springframework.data.domain.Pageable;

import java.util.List;

public interface MediaProvider {
    public List<MediaDto> searchMedia(Pageable pageable, SearchMediaRequest searchMediaRequest);

    public MediaDto getMedia(Long internalId);
}
