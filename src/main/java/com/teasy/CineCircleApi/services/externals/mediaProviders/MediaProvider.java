package com.teasy.CineCircleApi.services.externals.mediaProviders;

import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import org.springframework.data.domain.Pageable;

import java.util.List;

public interface MediaProvider {
    public List<MediaDto> searchMedia(Pageable pageable, MediaSearchRequest mediaSearchRequest);

    public MediaDto getMedia(Long internalId);
}
