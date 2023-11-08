package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.requests.LibraryAddMediaRequest;
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

import java.util.Optional;
import java.util.UUID;

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

    public void addToLibrary(UUID mediaId, LibraryAddMediaRequest libraryAddMediaRequest, String username) throws CustomException {
        // update existing library adding if exists
        var existingRecord = findExistingRecord(username, mediaId);
        if (existingRecord.isPresent()) {
            if (libraryAddMediaRequest != null) {
                var existingLibrary = existingRecord.get();
                existingLibrary.setComment(libraryAddMediaRequest.comment());
                existingLibrary.setRating(libraryAddMediaRequest.rating());
                libraryRepository.save(existingLibrary);
            }
        } else { // create new one if not
            if (libraryAddMediaRequest == null) {
                libraryAddMediaRequest = new LibraryAddMediaRequest(null, null);
            }
            libraryRepository.save(newLibrary(username, mediaId, libraryAddMediaRequest));
        }

    }

    public void removeFromLibrary(String username, UUID mediaId) {
        // get existing library record
        var existingRecord = findExistingRecord(username, mediaId);

        existingRecord.ifPresent(library -> libraryRepository.delete(library));
    }

    public Page<MediaShortDto> searchInLibrary(Pageable pageable,
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
        return records.map(library -> {
            var mediaDto = fromMediaEntityToMediaDto(library.getMedia());
            // add personal comment and rating
            mediaDto.setPersonalComment(library.getComment());
            mediaDto.setPersonalRating(library.getRating());
            return mediaDto;
        });
    }

    private Optional<Library> findExistingRecord(String username, UUID mediaId) {
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var matchingLibrary = newLibrary(username, mediaId, new LibraryAddMediaRequest(null, null));
        matchingLibrary.setAddedAt(null);
        return libraryRepository.findOne(Example.of(matchingLibrary, matcher));
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

    private Library newLibrary(String username, UUID mediaId, LibraryAddMediaRequest libraryAddMediaRequest) {
        var user = findUserByUsernameOrElseThrow(username);
        var media = findMediaByIdOrElseThrow(mediaId);
        return new Library(user, media, libraryAddMediaRequest.comment(), libraryAddMediaRequest.rating());
    }

    private User findUserByUsernameOrElseThrow(String username) {
        // check if user exist
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
    }

    private Media findMediaByIdOrElseThrow(UUID mediaId) {
        // check if media exists
        return mediaRepository
                .findById(mediaId)
                .orElseThrow(() -> CustomExceptionHandler.mediaWithIdNotFound(mediaId));
    }

    private MediaShortDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaShortDto.class);
    }
}
