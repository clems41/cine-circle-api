package com.teasy.CineCircleApi.services.externals.mediaProviders;

import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.dtos.requests.SearchMediaRequest;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

@Service
@ConfigurationProperties(prefix = "the-movie-db")
@Slf4j
public class TheMovieDb implements MediaProvider{
    private final static String baseUrl = "https://api.themoviedb.org/3";
    private final static String getMediaSuffix = "movie";
    private final static String searchMediaSuffix = "search/multi";
    private final static String searchMediaKey = "query";
    private final static String languageKey = "language";
    private final static String languageValue = "fr-FR";

    @Value("${the-movie-db.token}")
    private String token;
    MediaRepository mediaRepository;

    @Autowired
    public TheMovieDb(MediaRepository mediaRepository) {
        this.mediaRepository = mediaRepository;
    }

    @Override
    public Page<MediaShortDto> searchMedia(Pageable pageable, SearchMediaRequest searchMediaRequest) {
        var fullUrlPath = String.format("%s/%s?%s=%s&%s=%s",
                baseUrl,
                searchMediaSuffix,
                languageKey,
                languageValue,
                searchMediaKey,
                searchMediaRequest.query());
        HttpClient client = HttpClient.newHttpClient();
        HttpRequest request = HttpRequest.newBuilder()
                .GET()
                .uri(URI.create(fullUrlPath))
                .header(HttpHeaders.AUTHORIZATION, String.format("Bearer %s", token))
                .build();

        try {
            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
            var result = response.body();
            log.info(String.format("result = %s", result));
        } catch (Exception e) {
            throw new ResponseStatusException(
                    HttpStatus.INTERNAL_SERVER_ERROR,
                    String.format("cannot get response from provider for url %s : %s",
                            fullUrlPath,
                            e.getMessage()));
        }
        return null;
    }

    @Override
    public Media getMedia(Long internalId) throws ResponseStatusException {
        // prepare request
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setInternalId(internalId);

        // get media from to find external id
        var media = mediaRepository
                .findOne(Example.of(exampleMedia, matcher))
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("media with internalId %d cannot be found", internalId)));
        var externalId = media.getExternalId();

        // request media from provider API using externalId

        // define http client
        var fullUrlPath = String.format("%s/%s/%s?%s=%s",
                baseUrl,
                getMediaSuffix,
                media.getExternalId(),
                languageKey,
                languageValue);
            HttpClient client = HttpClient.newHttpClient();
            HttpRequest request = HttpRequest.newBuilder()
                    .GET()
                    .uri(URI.create(fullUrlPath))
                    .build();

            try {
                HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
                var result = response.body();
                log.info(String.format("result = %s", result));
            } catch (Exception e) {
                throw new ResponseStatusException(
                        HttpStatus.INTERNAL_SERVER_ERROR,
                        String.format("cannot get response from provider for url %s", fullUrlPath));
            }

        return null;
    }
}
