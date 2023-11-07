package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.CineCircleApiApplication;
import com.teasy.CineCircleApi.models.CustomPageImpl;
import com.teasy.CineCircleApi.models.dtos.MediaFullDto;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.models.dtos.requests.LibraryAddMediaRequest;
import com.teasy.CineCircleApi.models.dtos.responses.AuthSignInResponse;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.utils.Authenticator;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import com.teasy.CineCircleApi.utils.HttpUtils;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;
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
public class LibraryTest {
    @LocalServerPort
    private int port;

    @Autowired
    private TestRestTemplate restTemplate;

    @Autowired
    private MediaRepository mediaRepository;
    private final static String libraryUrl = "/library/";
    private final static String mediaUrl = "/medias/";
    private Authenticator authenticator;

    @BeforeEach
    public void setUp() {
        authenticator = new Authenticator(restTemplate, port);
    }

    @Test
    public void AddAndRemoveMultipleMedias() {
        /* Init */
        var dummyDataCreator = new DummyDataCreator(null, mediaRepository, null);
        var media1 = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media1 in database
        var media2 = dummyDataCreator.generateMedia(true, MediaTypeEnum.TV_SHOW); // create media2 in database
        var media3 = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media3 in database
        var nonExistingMediaId = UUID.randomUUID();

        /* Create user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Create Authorization header with JWT token */
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());

        /* Add media1 to library */
        ResponseEntity<String> addMedia1Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media1.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(new LibraryAddMediaRequest(null, null), headers),
                        String.class
                );
        Assertions.assertThat(addMedia1Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Add non-existing media to library */
        ResponseEntity<String> addNonExistingMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(nonExistingMediaId.toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(new LibraryAddMediaRequest(null, null), headers),
                        String.class
                );
        Assertions.assertThat(addNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Add media2 to library with comment and rating */
        var requestMedia2 = new LibraryAddMediaRequest(
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5));
        ResponseEntity<String> addMedia2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(requestMedia2, headers),
                        String.class
                );
        Assertions.assertThat(addMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse.getBody()).isNotNull();
        List<MediaShortDto> library = listLibraryResponse.getBody().stream().toList();
        Assertions.assertThat(library).hasSize(2);
        Assertions.assertThat(library.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that library contains media1
        // check that media2 is associated with comment and rating
        var media2FromList = library.stream()
                .filter(mediaDto -> Objects.equals(mediaDto.getId(), media2.getId().toString()))
                .findAny();
        Assertions.assertThat(media2FromList.isPresent()).isTrue(); // check that library contains media2
        Assertions.assertThat(media2FromList.get().getPersonalComment()).isEqualTo(requestMedia2.comment());
        Assertions.assertThat(media2FromList.get().getPersonalRating()).isEqualTo(requestMedia2.rating());

        /* Remove non-existing media from library */
        ResponseEntity<String> removeNonExistingMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(nonExistingMediaId.toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Remove media2 from library */
        ResponseEntity<String> removeMedia2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Remove media2 a 2nd time from library */
        ResponseEntity<String> removeMedia2SecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2SecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse2 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse2.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse2.getBody()).isNotNull();
        List<MediaShortDto> library2 = listLibraryResponse2.getBody().stream().toList();
        Assertions.assertThat(library2).hasSize(1);
        Assertions.assertThat(library2.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that library contains media1

        /* Add media3 to library with comment and rating */
        var requestMedia3 = new LibraryAddMediaRequest(
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5));
        ResponseEntity<String> addMedia3Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media3.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(requestMedia3, headers),
                        String.class
                );
        Assertions.assertThat(addMedia3Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse3 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse3.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse3.getBody()).isNotNull();
        List<MediaShortDto> library3 = listLibraryResponse3.getBody().stream().toList();
        Assertions.assertThat(library3).hasSize(2);
        Assertions.assertThat(library3.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that library contains media1
        // check that media3 is associated with comment and rating
        var media3FromList = library3.stream()
                .filter(mediaDto -> Objects.equals(mediaDto.getId(), media3.getId().toString()))
                .findAny();
        Assertions.assertThat(media3FromList.isPresent()).isTrue(); // check that library contains media2
        Assertions.assertThat(media3FromList.get().getPersonalComment()).isEqualTo(requestMedia3.comment());
        Assertions.assertThat(media3FromList.get().getPersonalRating()).isEqualTo(requestMedia3.rating());

        /* Add media3 to library a second time with new comment and new rating */
        var requestMedia3SecondTime = new LibraryAddMediaRequest(
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5));
        ResponseEntity<String> addMedia3SecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media3.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(requestMedia3SecondTime, headers),
                        String.class
                );
        Assertions.assertThat(addMedia3SecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibrarySecondTimeResponse3 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibrarySecondTimeResponse3.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibrarySecondTimeResponse3.getBody()).isNotNull();
        List<MediaShortDto> library3SecondTIme = listLibrarySecondTimeResponse3.getBody().stream().toList();
        Assertions.assertThat(library3SecondTIme).hasSize(2);
        Assertions.assertThat(library3SecondTIme.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that library contains media1
        // check that media3 is associated with comment and rating
        var media3FromListSecondTime = library3SecondTIme.stream()
                .filter(mediaDto -> Objects.equals(mediaDto.getId(), media3.getId().toString()))
                .findAny();
        Assertions.assertThat(media3FromListSecondTime.isPresent()).isTrue(); // check that library contains media2
        Assertions.assertThat(media3FromListSecondTime.get().getPersonalComment()).isEqualTo(requestMedia3SecondTime.comment());
        Assertions.assertThat(media3FromListSecondTime.get().getPersonalRating()).isEqualTo(requestMedia3SecondTime.rating());
    }
    @Test
    public void AddMedia_ShouldContainsCommentAndRatingWhenGettingMedia() {
        /* Init */
        var dummyDataCreator = new DummyDataCreator(null, mediaRepository, null);
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media in database

        /* Create user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Create Authorization header with JWT token */
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());

        /* Add media to library */
        var addMediaRequest = new LibraryAddMediaRequest(
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5));
        ResponseEntity<String> addMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(addMediaRequest, headers),
                        String.class
                );
        Assertions.assertThat(addMediaResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get media and check rating and comment */
        ResponseEntity<MediaFullDto> getMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(mediaUrl).concat(media.getId().toString()),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        MediaFullDto.class
                );
        Assertions.assertThat(getMediaResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(getMediaResponse.getBody()).isNotNull();
        Assertions.assertThat(getMediaResponse.getBody().getPersonalComment()).isEqualTo(addMediaRequest.comment());
        Assertions.assertThat(getMediaResponse.getBody().getPersonalRating()).isEqualTo(addMediaRequest.rating());

        /* Add same media to library with new comment and new rating */
        var addMediaSecondTimeRequest = new LibraryAddMediaRequest(
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5));
        ResponseEntity<String> addMediaSecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(libraryUrl).concat(media.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(addMediaSecondTimeRequest, headers),
                        String.class
                );
        Assertions.assertThat(addMediaSecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get media and check new rating and new comment */
        ResponseEntity<MediaFullDto> getMediaSecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(mediaUrl).concat(media.getId().toString()),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        MediaFullDto.class
                );
        Assertions.assertThat(getMediaSecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(getMediaSecondTimeResponse.getBody()).isNotNull();
        Assertions.assertThat(getMediaSecondTimeResponse.getBody().getPersonalComment()).isEqualTo(addMediaSecondTimeRequest.comment());
        Assertions.assertThat(getMediaSecondTimeResponse.getBody().getPersonalRating()).isEqualTo(addMediaSecondTimeRequest.rating());
    }
}