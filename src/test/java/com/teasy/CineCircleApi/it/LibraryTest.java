package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.CineCircleApiApplication;
import com.teasy.CineCircleApi.models.CustomPageImpl;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.models.dtos.responses.AuthSignInResponse;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import org.apache.commons.lang3.RandomStringUtils;
import org.assertj.core.api.Assertions;
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

    private final static String authSignUpUrl = "/auth/sign-up";
    private final static String authSignInUrl = "/auth/sign-in";
    private final static String libraryUrl = "/library/";

    @Test
    public void AddAndRemoveMultipleMedias() {
        /* Init */
        var dummyDataCreator = new DummyDataCreator(null, mediaRepository, null);
        var media1 = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media1 in database
        var media2 = dummyDataCreator.generateMedia(true, MediaTypeEnum.TV_SHOW); // create media2 in database
        var media3 = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media3 in database
        var nonExistingMediaId = UUID.randomUUID();

        /* Create user */
        var authSignUpRequest = generateAuthSignUpRequest();
        ResponseEntity<UserFullInfoDto> signUpResponse = this.restTemplate
                .postForEntity(
                        getTestingUrl().concat(authSignUpUrl),
                        authSignUpRequest,
                        UserFullInfoDto.class
                );
        Assertions.assertThat(signUpResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(signUpResponse.getBody()).isNotNull();
        Assertions.assertThat(signUpResponse.getBody().getUsername()).isEqualTo(authSignUpRequest.username());
        Assertions.assertThat(signUpResponse.getBody().getEmail()).isEqualTo(authSignUpRequest.email());
        Assertions.assertThat(signUpResponse.getBody().getDisplayName()).isEqualTo(authSignUpRequest.displayName());
        Assertions.assertThat(signUpResponse.getBody().getTopicName()).isNotEmpty();
        Assertions.assertThat(signUpResponse.getBody().getId()).isNotEmpty();

        /* Get token for new created user */
        ResponseEntity<AuthSignInResponse> signInResponse = this.restTemplate
                .withBasicAuth(authSignUpRequest.username(), authSignUpRequest.password())
                .getForEntity(
                        getTestingUrl().concat(authSignInUrl),
                        AuthSignInResponse.class
                );
        Assertions.assertThat(signInResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(signInResponse.getBody()).isNotNull();
        var jwtToken = signInResponse.getBody().getToken().tokenString();
        Assertions.assertThat(jwtToken).isNotEmpty();

        /* Create Authorization header with JWT token */
        HttpHeaders headers = new HttpHeaders();
        headers.setBearerAuth(jwtToken);

        /* Add media1 to library */
        ResponseEntity<String> addMedia1Response = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl).concat(media1.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMedia1Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Add non-existing media to library */
        ResponseEntity<String> addNonExistingMediaResponse = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl).concat(nonExistingMediaId.toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Add media2 to library */
        ResponseEntity<String> addMedia2Response = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaDto>> listLibraryResponse = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse.getBody()).isNotNull();
        List<MediaDto> library = listLibraryResponse.getBody().stream().toList();
        Assertions.assertThat(library).hasSize(2);
        Assertions.assertThat(library.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
                )).isTrue(); // check that library contains media1
        Assertions.assertThat(library.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media2.getId().toString())
                )).isTrue(); // check that library contains media2

        /* Remove non-existing media from library */
        ResponseEntity<String> removeNonExistingMediaResponse = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl).concat(nonExistingMediaId.toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Remove media2 from library */
        ResponseEntity<String> removeMedia2Response = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Remove media2 a 2nd time from library */
        ResponseEntity<String> removeMedia2SecondTimeResponse = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2SecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaDto>> listLibraryResponse2 = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse2.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse2.getBody()).isNotNull();
        List<MediaDto> library2 = listLibraryResponse2.getBody().stream().toList();
        Assertions.assertThat(library2).hasSize(1);
        Assertions.assertThat(library2.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that library contains media1

        /* Add media3 to library */
        ResponseEntity<String> addMedia3Response = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl).concat(media3.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMedia3Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaDto>> listLibraryResponse3 = this.restTemplate
                .exchange(
                        getTestingUrl().concat(libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse3.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse3.getBody()).isNotNull();
        List<MediaDto> library3 = listLibraryResponse3.getBody().stream().toList();
        Assertions.assertThat(library3).hasSize(2);
        Assertions.assertThat(library3.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media1.getId().toString())
        )).isTrue(); // check that library contains media1
        Assertions.assertThat(library3.stream().anyMatch(
                mediaDto -> Objects.equals(mediaDto.getId(), media3.getId().toString())
        )).isTrue(); // check that library contains media3
    }

    private AuthSignUpRequest generateAuthSignUpRequest() {
        var username = RandomStringUtils.random(10, true, true);
        var email = String.format("%s@%s.com",
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(5, true, false));
        var password = RandomStringUtils.random(15, true, true);
        var displayName = RandomStringUtils.random(20, true, true);
        return new AuthSignUpRequest(username, email, password, displayName);
    }

    private String getTestingUrl() {
        return "http://localhost:".concat(String.valueOf(port)).concat("/api/v1");
    }
}
