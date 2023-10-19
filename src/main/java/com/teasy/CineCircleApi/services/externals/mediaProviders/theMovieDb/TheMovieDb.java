package com.teasy.CineCircleApi.services.externals.mediaProviders.theMovieDb;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.dtos.requests.SearchMediaRequest;
import com.teasy.CineCircleApi.models.enums.MediaType;
import com.teasy.CineCircleApi.models.externals.theMovieDb.Genre;
import com.teasy.CineCircleApi.models.externals.theMovieDb.GenreMovieListResponse;
import com.teasy.CineCircleApi.models.externals.theMovieDb.SearchResponse;
import com.teasy.CineCircleApi.models.externals.theMovieDb.SearchResult;
import com.teasy.CineCircleApi.models.utils.CustomHttpClientSendRequest;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.utils.CustomHttpClient;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.io.IOException;
import java.io.InputStream;
import java.util.*;

@Service
@Slf4j
public class TheMovieDb implements MediaProvider {
    private final static String stringArrayDelimiter = ",";
    private final static String baseUrl = "https://api.themoviedb.org/3";
    private final static String getMediaSuffix = "movie";
    private final static String searchMediaSuffix = "search/multi";
    private final static String searchMediaKey = "query";
    private final static String languageKey = "language";
    private final static String languageValue = "fr-FR";
    private final static String tokenPrefix = "Bearer";
    private final static String imageUrlPrefix = "https://image.tmdb.org/t/p/w500";

    private List<Genre> genres;

    @Value("${the-movie-db.genres-input-file-path}")
    private String genresInputFilePath;

    @Value("${the-movie-db.token}")
    private String token;
    MediaRepository mediaRepository;
    CustomHttpClient customHttpClient;

    @Autowired
    public TheMovieDb(MediaRepository mediaRepository,
                      CustomHttpClient customHttpClient) {
        this.mediaRepository = mediaRepository;
        this.customHttpClient = customHttpClient;
    }

    @Override
    public List<MediaDto> searchMedia(Pageable pageable, SearchMediaRequest searchMediaRequest) {
        // define request
        var url = String.format("%s/%s", baseUrl, searchMediaSuffix);
        var queryParameters = getDefaultQueryParameters();
        queryParameters.put(searchMediaKey, searchMediaRequest.query());
        var request = new CustomHttpClientSendRequest(HttpMethod.GET, url, queryParameters, getDefaultAuthorizationHeader());

        // send request
        SearchResponse response = customHttpClient.sendRequest(request, SearchResponse.class);
        if (response.getResults() == null) {
            return new ArrayList<>();
        }

        // get genres list before mapping result
        setGenresFromFile();

        List<MediaDto> result = new ArrayList<>();
        response.getResults().forEach(searchResult -> {
            // store result in database with internalId if not already exists
            var existingMedia = findMediaFromSearchResult(searchResult);
            if (existingMedia.isEmpty()) {
                 var newMedia = mediaRepository.save(fromSearchResultToMediaEntity(searchResult));
                 result.add(fromMediaEntityToMediaDto(newMedia));
            } else {
                result.add(fromMediaEntityToMediaDto(existingMedia.get()));
            }

        });
        return result;
    }

    @Override
    public MediaDto getMedia(Long id) throws ResponseStatusException {
        // build example matcher with id
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setId(id);

        // get media from database
        var media = mediaRepository
                .findOne(Example.of(exampleMedia, matcher))
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("media with id %d cannot be found", id)));
        return fromMediaEntityToMediaDto(media);
    }

    private void setGenresFromFile() {
        if (genres != null) {
            return;
        }
        try(InputStream in = Thread.currentThread()
                .getContextClassLoader()
                .getResourceAsStream(genresInputFilePath)) {
            ObjectMapper mapper = new ObjectMapper()
                    .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
            GenreMovieListResponse result = mapper.readValue(in, GenreMovieListResponse.class);
            genres = result.getGenres();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private Optional<Media> findMediaFromSearchResult(SearchResult searchResult) {
        // build example matcher with external id and media provider
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setExternalId(String.valueOf(searchResult.getId()));
        exampleMedia.setMediaProvider(com.teasy.CineCircleApi.models.enums.MediaProvider.THE_MOVIE_DATABASE.name());

        return mediaRepository
                .findOne(Example.of(exampleMedia, matcher));
    }

    private Media fromSearchResultToMediaEntity(SearchResult searchResult) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        Media media = mapper.convertValue(searchResult, Media.class);
        switch (searchResult.getMediaType()) {
            case "tv" -> media.setMediaType(MediaType.TV_SHOW.name());
            case "movie" -> media.setMediaType(MediaType.MOVIE.name());
        }
        media.setOriginalTitle(searchResult.getOriginalTitle() != null ? searchResult.getOriginalTitle() : searchResult.getOriginalName());
        media.setTitle(searchResult.getTitle() != null ? searchResult.getTitle() : searchResult.getName());
        media.setExternalId(String.valueOf(searchResult.getId()));
        media.setMediaProvider(com.teasy.CineCircleApi.models.enums.MediaProvider.THE_MOVIE_DATABASE.name());
        media.setPosterUrl(getCompleteImageUrl(searchResult.getPosterPath()));
        media.setBackdropUrl(getCompleteImageUrl(searchResult.getBackdropPath()));
        media.setReleaseDate(searchResult.getReleaseDate() != null ? searchResult.getReleaseDate() : searchResult.getFirstAirDate());
        media.setOriginalLanguage(searchResult.getOriginalLanguage());
        media.setVoteAverage(searchResult.getVoteAverage());
        media.setVoteCount(searchResult.getVoteCount());
        if (searchResult.getOriginCountry() != null) {
            media.setOriginCountry(String.join(stringArrayDelimiter, searchResult.getOriginCountry()));
        }
        if (searchResult.getGenreIds() != null) {
            media.setGenres(String.join(stringArrayDelimiter,
                    searchResult
                            .getGenreIds()
                            .stream()
                            .map(this::getGenreFromId)
                            .filter(s -> !s.isEmpty())
                            .toList()
            ));
        }
        return media;
    }

    private String getGenreFromId(Long genreId) {
        return genres.stream()
                .filter(genre -> Objects.equals(genre.getId(), genreId))
                .findAny()
                .map(Genre::getName)
                .orElse("");
    }

    private String getCompleteImageUrl(String posterUrl) {
        return String.format("%s%s", imageUrlPrefix, posterUrl);
    }

    private MediaDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        return mapper.convertValue(media, MediaDto.class);
    }

    private Map<String, String> getDefaultQueryParameters() {
        Map<String, String> queryParameters = new HashMap<>();
        queryParameters.put(languageKey, languageValue);
        return queryParameters;
    }

    private String getDefaultAuthorizationHeader() {
        return String.format("%s %s", tokenPrefix, token);
    }
}
