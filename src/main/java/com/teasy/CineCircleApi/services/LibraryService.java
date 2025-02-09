package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.requests.LibraryAddMediaRequest;
import com.teasy.CineCircleApi.models.dtos.requests.LibrarySearchRequest;
import com.teasy.CineCircleApi.models.entities.Library;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.LibraryRepository;
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
    private final LibraryRepository libraryRepository;
    private final UserService userService;
    private final MediaService mediaService;

    @Autowired
    public LibraryService(LibraryRepository libraryRepository, UserService userService, MediaService mediaService) {
        this.libraryRepository = libraryRepository;
        this.userService = userService;
        this.mediaService = mediaService;
    }

    public void addToLibrary(UUID mediaId, LibraryAddMediaRequest libraryAddMediaRequest, String username) throws ExpectedException {
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

    public void removeFromLibrary(String username, UUID mediaId) throws ExpectedException {
        // get existing library record
        var existingRecord = findExistingRecord(username, mediaId);
        existingRecord.ifPresent(libraryRepository::delete);
    }

    public Page<MediaShortDto> searchInLibrary(Pageable pageable,
                                               LibrarySearchRequest librarySearchRequest,
                                               String username) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);

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

    public boolean isInLibrary(String username, UUID mediaId) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        return libraryRepository.existsByUser_IdAndMedia_Id(user.getId(), mediaId);
    }

    private Optional<Library> findExistingRecord(String username, UUID mediaId) throws ExpectedException {
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

    private Library newLibrary(String username, UUID mediaId, LibraryAddMediaRequest libraryAddMediaRequest) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        var media = mediaService.findMediaByIdOrElseThrow(mediaId);
        return new Library(user, media, libraryAddMediaRequest.comment(), libraryAddMediaRequest.rating());
    }

    private MediaShortDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(media, MediaShortDto.class);
    }
}
