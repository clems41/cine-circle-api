package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.LibrarySearchRequest;
import com.teasy.CineCircleApi.models.entities.Library;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.repositories.LibraryRepository;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

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

    public void addToLibrary(String username, Long mediaId) throws CustomException {
        libraryRepository.save(newLibrary(username, mediaId));
    }

    public void removeFromLibrary(String username, Long mediaId) {
        // get existing library record
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingLibrary = newLibrary(username, mediaId);
        matchingLibrary.setAddedAt(null);
        var existingRecord = libraryRepository.findOne(Example.of(matchingLibrary, matcher));

        existingRecord.ifPresent(library -> libraryRepository.delete(library));
    }

    public Page<MediaDto> searchInLibrary(Pageable pageable,
                                          LibrarySearchRequest librarySearchRequest,
                                          String username) {
        var user = findUserByUsernameOrElseThrow(username);

        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues()
                .withIgnoreCase()
                .withStringMatcher(ExampleMatcher.StringMatcher.CONTAINING);

        // creating matching library based on request
        var matchingLibrary = new Library();
        matchingLibrary.setUser(user);
        matchingLibrary.setMedia(createMatchingMedia(librarySearchRequest));

        var records = libraryRepository.findAll(Example.of(matchingLibrary, matcher), pageable);
        return records.map(library -> fromMediaEntityToMediaDto(library.getMedia()));
    }

    private Media createMatchingMedia(LibrarySearchRequest librarySearchRequest) {
        var matchingMedia = new Media();
        if (librarySearchRequest.query() != null && !librarySearchRequest.query().isEmpty()) {
            matchingMedia.setTitle(librarySearchRequest.query());
        } else if (librarySearchRequest.genre() != null && !librarySearchRequest.genre().isEmpty()) {
            matchingMedia.setGenres(librarySearchRequest.genre());
        }
        return matchingMedia;
    }

    private Library newLibrary(String username, Long mediaId) {
        var user = findUserByUsernameOrElseThrow(username);
        var media = findMediaByIdOrElseThrow(mediaId);
        return new Library(user, media);
    }

    private User findUserByUsernameOrElseThrow(String username) {
        // check if user exist
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
    }

    private Media findMediaByIdOrElseThrow(Long mediaId) {
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
