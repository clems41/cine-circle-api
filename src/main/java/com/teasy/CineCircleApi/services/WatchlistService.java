package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.entities.Watchlist;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.repositories.WatchlistRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

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

    public void addToWatchlist(String username, Long mediaId) throws ResponseStatusException {
        // check if user exist
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("user with username %s cannot be found", username)));

        // check if media exists
        var media = mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("media with id %d cannot be found", mediaId)));

        // add to watchlist if both exists
        var watchlist = new Watchlist(user, media);
        watchlistRepository.save(watchlist);
    }

    public void removeFromWatchlist(String username, Long mediaId) {
        // check if user exist
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("user with username %s cannot be found", username)));

        // check if media exists
        var media = mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("media with id %d cannot be found", mediaId)));

        // get existing watchlist record
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingUser = new User();
        matchingUser.setId(user.getId());
        var matchingMedia = new Media();
        matchingMedia.setId(media.getId());
        var matchingWatchlist = new Watchlist(matchingUser, matchingMedia);
        matchingWatchlist.setAddedAt(null);
        var existingRecord = watchlistRepository.findOne(Example.of(matchingWatchlist, matcher));

        existingRecord.ifPresent(watchlist -> watchlistRepository.delete(watchlist));
    }

    public Page<MediaDto> getWatchlist(Pageable pageable, String username) {
        // check if user exist
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("user with username %s cannot be found", username)));


        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingWatchlist = new Watchlist();
        matchingWatchlist.setUser(user);

        var records = watchlistRepository.findAll(Example.of(matchingWatchlist, matcher), pageable);
        return records.map(watchlist -> fromMediaEntityToMediaDto(watchlist.getMedia()));
    }

    private MediaDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaDto.class);
    }
}
