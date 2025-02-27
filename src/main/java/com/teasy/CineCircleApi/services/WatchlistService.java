package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Watchlist;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.WatchlistRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class WatchlistService {

    private final WatchlistRepository watchlistRepository;
    private final MediaService mediaService;
    private final UserService userService;

    @Autowired
    public WatchlistService(WatchlistRepository watchlistRepository, MediaService mediaService, UserService userService) {
        this.watchlistRepository = watchlistRepository;
        this.mediaService = mediaService;
        this.userService = userService;
    }

    public void addToWatchlist(String username, UUID mediaId) throws ExpectedException {
        watchlistRepository.save(newWatchlist(username, mediaId));
    }

    public void removeFromWatchlist(String username, UUID mediaId) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        mediaService.findMediaByIdOrElseThrow(mediaId); // check if media exists
        var existingRecord = watchlistRepository.findByUser_IdAndMedia_Id(user.getId(), mediaId);
        existingRecord.ifPresent(watchlistRepository::delete);
    }

    public Page<MediaShortDto> getWatchlist(Pageable pageable, String username) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);

        var records = watchlistRepository.findByUser_Id(user.getId(), pageable);
        return records.map(watchlist -> fromMediaEntityToMediaDto(watchlist.getMedia()));
    }

    public boolean isInWatchlist(String username, UUID mediaId) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        return watchlistRepository.existsByUser_IdAndMedia_Id(user.getId(), mediaId);
    }

    private Watchlist newWatchlist(String username, UUID mediaId) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        var media = mediaService.findMediaByIdOrElseThrow(mediaId);
        return new Watchlist(user, media);
    }

    private MediaShortDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaShortDto.class);
    }
}
