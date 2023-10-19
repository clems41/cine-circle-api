package com.teasy.CineCircleApi.services.externals.mediaProviders;

import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.dtos.requests.SearchMediaRequest;
import com.teasy.CineCircleApi.models.entities.MediaType;
import com.teasy.CineCircleApi.models.externals.theMovieDb.SearchResponse;
import com.teasy.CineCircleApi.models.externals.theMovieDb.SearchResult;
import com.teasy.CineCircleApi.models.utils.CustomHttpClientSendRequest;
import com.teasy.CineCircleApi.repositories.MediaRepository;
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

import java.util.*;

@Service
@ConfigurationProperties(prefix = "the-movie-db")
@Slf4j
public class TheMovieDb implements MediaProvider {
    private final static String baseUrl = "https://api.themoviedb.org/3";
    private final static String getMediaSuffix = "movie";
    private final static String searchMediaSuffix = "search/multi";
    private final static String searchMediaKey = "query";
    private final static String languageKey = "language";
    private final static String languageValue = "fr-FR";
    private final static String tokenPrefix = "Bearer";
    private final static String imageUrlPrefix = "https://image.tmdb.org/t/p/w500";

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

        List<MediaDto> result = new ArrayList<>();
        response.getResults().forEach(searchResult -> {
            // store result in database with internalId if not already exists
            var existingMedia = findMediaFromSearchResult(searchResult);
            if (existingMedia.isEmpty()) {
                 var newMedia = mediaRepository.save(fromSearchResultToMediaEntity(searchResult));
                 result.add(fromMediaEntityToMediaShortDto(newMedia));
            } else {
                result.add(fromMediaEntityToMediaShortDto(existingMedia.get()));
            }

        });
        return result;
    }

    @Override
    public MediaDto getMedia(Long internalId) throws ResponseStatusException {
        // build example matcher with id
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setInternalId(internalId);

        // get media from database
        var media = mediaRepository
                .findOne(Example.of(exampleMedia, matcher))
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("media with internalId %d cannot be found", internalId)));
        return fromMediaEntityToMediaShortDto(media);
    }

    private Optional<Media> findMediaFromSearchResult(SearchResult searchResult) {
        // build example matcher with external id and media provider
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setExternalId(String.valueOf(searchResult.getId()));
        exampleMedia.setMediaProvider(com.teasy.CineCircleApi.models.entities.MediaProvider.THE_MOVIE_DATABASE);

        return mediaRepository
                .findOne(Example.of(exampleMedia, matcher));
    }

    private Media fromSearchResultToMediaEntity(SearchResult searchResult) {
        var media = new Media();
        switch (searchResult.getMediaType()) {
            case "tv" -> media.setMediaType(MediaType.TV_SHOW);
            case "movie" -> media.setMediaType(MediaType.MOVIE);
        }
        media.setSynopsis(searchResult.getOverview());
        media.setTitle(searchResult.getTitle());
        media.setOriginalTitle(searchResult.getOriginalTitle());
        media.setExternalId(String.valueOf(searchResult.getId()));
        media.setReleaseDate(searchResult.getReleaseDate());
        media.setMediaProvider(com.teasy.CineCircleApi.models.entities.MediaProvider.THE_MOVIE_DATABASE);
        media.setPosterUrl(getCompletePosterUrl(searchResult.getPosterPath()));
        return media;
    }

    private String getCompletePosterUrl(String posterUrl) {
        return String.format("%s%s", imageUrlPrefix, posterUrl);
    }

    private MediaDto fromMediaEntityToMediaShortDto(Media media) {
        return new MediaDto(
                media.getInternalId(),
                media.getTitle(),
                media.getOriginalTitle(),
                media.getPosterUrl()
        );
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
