package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.entities.Watchlist;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.repositories.WatchlistRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

@Service
public class WatchlistService {

    WatchlistRepository watchlistRepository;
    MediaRepository mediaRepository;
    UserRepository userRepository;

    @Autowired
    public WatchlistService(WatchlistRepository watchlistRepository, MediaRepository mediaRepository, UserRepository userRepository) {
        this.watchlistRepository = watchlistRepository;
        this.mediaRepository = mediaRepository;
        this.userRepository = userRepository;
    }

    public void addToWatchlist(String username, Long mediaId) throws CustomException {
        watchlistRepository.save(newWatchlist(username, mediaId));
    }

    public void removeFromWatchlist(String username, Long mediaId) {
        // get existing watchlist record
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingWatchlist = newWatchlist(username, mediaId);
        matchingWatchlist.setAddedAt(null);
        var existingRecord = watchlistRepository.findOne(Example.of(matchingWatchlist, matcher));

        existingRecord.ifPresent(watchlist -> watchlistRepository.delete(watchlist));
    }

    public Page<MediaDto> getWatchlist(Pageable pageable, String username) {
        var user = getUserWithUsernameOrElseThrow(username);

        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingWatchlist = new Watchlist();
        matchingWatchlist.setUser(user);

        var records = watchlistRepository.findAll(Example.of(matchingWatchlist, matcher), pageable);
        return records.map(watchlist -> fromMediaEntityToMediaDto(watchlist.getMedia()));
    }

    private Watchlist newWatchlist(String username, Long mediaId) {
        var user = getUserWithUsernameOrElseThrow(username);
        var media = getMediaWithIdOrElseThrow(mediaId);
        return new Watchlist(user, media);
    }

    private User getUserWithUsernameOrElseThrow(String username) {
        // check if user exist
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
    }

    private Media getMediaWithIdOrElseThrow(Long mediaId) {
        // check if media exists
        return mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> CustomExceptionHandler.mediaWithIdNotFound(mediaId));
    }

    private MediaDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaDto.class);
    }
}
