package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.CineCircleApiApplication;
import com.teasy.CineCircleApi.models.CustomPageImpl;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.utils.Authenticator;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import com.teasy.CineCircleApi.utils.HttpUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.*;
import org.springframework.test.context.ActiveProfiles;

import java.util.List;
import java.util.Objects;
import java.util.UUID;

@ActiveProfiles("test")
@SpringBootTest(classes = CineCircleApiApplication.class, webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public class WatchlistTest {
    @LocalServerPort
    private int port;

    @Autowired
    private TestRestTemplate restTemplate;

    @Autowired
    private MediaRepository mediaRepository;
    private Authenticator authenticator;
    private DummyDataCreator dummyDataCreator;

    @BeforeEach
    public void setUp() {
        authenticator = new Authenticator(restTemplate, port);
        dummyDataCreator = new DummyDataCreator(null, mediaRepository, null, null);
    }

    @Test
    public void AddAndRemoveMultipleMedias() {
        /* Init */
        var media1 = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media1 in database
        var media2 = dummyDataCreator.generateMedia(true, MediaTypeEnum.TV_SHOW); // create media2 in database
        var media3 = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media3 in database
        var nonExistingMediaId = UUID.randomUUID();

        /* Create user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Create Authorization header with JWT token */
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());

        /* Add media1 to watchlist */
        ResponseEntity<String> addMedia1Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat(media1.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMedia1Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Add non-existing media to watchlist */
        ResponseEntity<String> addNonExistingMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat(nonExistingMediaId.toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Add media2 to watchlist */
        ResponseEntity<String> addMedia2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat(media2.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from watchlist */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listWatchlistResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listWatchlistResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listWatchlistResponse.getBody()).isNotNull();
        List<MediaShortDto> watchlist = listWatchlistResponse.getBody().stream().toList();
        Assertions.assertThat(watchlist).hasSize(2);
        Assertions.assertThat(watchlist.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
                )).isTrue(); // check that watchlist contains media1
        Assertions.assertThat(watchlist.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media2.getId().toString())
                )).isTrue(); // check that watchlist contains media2

        /* Remove non-existing media from watchlist */
        ResponseEntity<String> removeNonExistingMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat(nonExistingMediaId.toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Remove media2 from watchlist */
        ResponseEntity<String> removeMedia2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Remove media2 a 2nd time from watchlist */
        ResponseEntity<String> removeMedia2SecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2SecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from watchlist */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listWatchlistResponse2 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listWatchlistResponse2.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listWatchlistResponse2.getBody()).isNotNull();
        List<MediaShortDto> watchlist2 = listWatchlistResponse2.getBody().stream().toList();
        Assertions.assertThat(watchlist2).hasSize(1);
        Assertions.assertThat(watchlist2.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that watchlist2 contains media1

        /* Add media3 to watchlist */
        ResponseEntity<String> addMedia3Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat(media3.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMedia3Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from watchlist */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listWatchlistResponse3 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listWatchlistResponse3.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listWatchlistResponse3.getBody()).isNotNull();
        List<MediaShortDto> watchlist3 = listWatchlistResponse3.getBody().stream().toList();
        Assertions.assertThat(watchlist3).hasSize(2);
        Assertions.assertThat(watchlist3.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that watchlist3 contains media1
        Assertions.assertThat(watchlist3.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media3.getId().toString())
        )).isTrue(); // check that watchlist3 contains media3
    }
}
