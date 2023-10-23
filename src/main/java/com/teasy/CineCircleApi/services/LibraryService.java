package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.entities.Library;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.entities.Watchlist;
import com.teasy.CineCircleApi.repositories.LibraryRepository;
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
public class LibraryService {

    LibraryRepository libraryRepository;
    MediaRepository mediaRepository;
    UserRepository userRepository;

    @Autowired
    public LibraryService(LibraryRepository libraryRepository, MediaRepository mediaRepository, UserRepository userRepository) {
        this.libraryRepository = libraryRepository;
        this.mediaRepository = mediaRepository;
        this.userRepository = userRepository;
    }

    public void addToLibrary(String username, Long mediaId) throws ResponseStatusException {
        var user = getUserAndThrownIfNotExist(username);
        var media = getMediaAndThrownIfNotExist(mediaId);

        // add to watchlist if both exists
        var libraryRecord = new Library(user, media);
        libraryRepository.save(libraryRecord);
    }

    public void removeFromLibrary(String username, Long mediaId) {
        var user = getUserAndThrownIfNotExist(username);
        var media = getMediaAndThrownIfNotExist(mediaId);

        // get existing library record
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingUser = new User();
        matchingUser.setId(user.getId());
        var matchingMedia = new Media();
        matchingMedia.setId(media.getId());
        var matchingLibrary = new Library(matchingUser, matchingMedia);
        matchingLibrary.setAddedAt(null);
        var existingRecord = libraryRepository.findOne(Example.of(matchingLibrary, matcher));

        existingRecord.ifPresent(library -> libraryRepository.delete(library));
    }

    public Page<MediaDto> getLibrary(Pageable pageable, String username) {
        var user = getUserAndThrownIfNotExist(username);

        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingLibrary = new Library();
        matchingLibrary.setUser(user);

        var records = libraryRepository.findAll(Example.of(matchingLibrary, matcher), pageable);
        return records.map(library -> fromMediaEntityToMediaDto(library.getMedia()));
    }

    private User getUserAndThrownIfNotExist(String username) {
        // check if user exist
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("user with username %s cannot be found", username)));
    }

    private Media getMediaAndThrownIfNotExist(Long mediaId) {
        // check if media exists
        return mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("media with id %d cannot be found", mediaId)));
    }

    private MediaDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaDto.class);
    }
}
