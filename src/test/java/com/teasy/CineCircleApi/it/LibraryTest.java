package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.MediaFullDto;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.LibraryAddMediaRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.utils.CustomPageImpl;
import com.teasy.CineCircleApi.utils.HttpUtils;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import java.util.UUID;

public class LibraryTest extends IntegrationTestAbstract {
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

        /* Add media1 to library */
        ResponseEntity<String> addMedia1Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media1.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(new LibraryAddMediaRequest(null, null), headers),
                        String.class
                );
        Assertions.assertThat(addMedia1Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Add non-existing media to library */
        ResponseEntity<String> addNonExistingMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(nonExistingMediaId.toString()),
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
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(requestMedia2, headers),
                        String.class
                );
        Assertions.assertThat(addMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl),
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

        /* Remove non-existing media from library */
        ResponseEntity<String> removeNonExistingMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(nonExistingMediaId.toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Remove media2 from library */
        ResponseEntity<String> removeMedia2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Remove media2 a 2nd time from library */
        ResponseEntity<String> removeMedia2SecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2SecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse2 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl),
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
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media3.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(requestMedia3, headers),
                        String.class
                );
        Assertions.assertThat(addMedia3Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse3 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl),
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

        /* Add media3 to library a second time with new comment and new rating */
        var requestMedia3SecondTime = new LibraryAddMediaRequest(
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5));
        ResponseEntity<String> addMedia3SecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media3.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(requestMedia3SecondTime, headers),
                        String.class
                );
        Assertions.assertThat(addMedia3SecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibrarySecondTimeResponse3 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl),
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
    }
    @Test
    public void AddMedia_ShouldContainsCommentAndRatingWhenGettingMedia() {
        /* Init */
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
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(addMediaRequest, headers),
                        String.class
                );
        Assertions.assertThat(addMediaResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get media and check rating and comment */
        ResponseEntity<MediaFullDto> getMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.mediaUrl).concat(media.getId().toString()),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        MediaFullDto.class
                );
        Assertions.assertThat(getMediaResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(getMediaResponse.getBody()).isNotNull();

        /* Add same media to library with new comment and new rating */
        var addMediaSecondTimeRequest = new LibraryAddMediaRequest(
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5));
        ResponseEntity<String> addMediaSecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl).concat(media.getId().toString()),
                        HttpMethod.POST,
                        new HttpEntity<>(addMediaSecondTimeRequest, headers),
                        String.class
                );
        Assertions.assertThat(addMediaSecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get media and check new rating and new comment */
        ResponseEntity<MediaFullDto> getMediaSecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.mediaUrl).concat(media.getId().toString()),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        MediaFullDto.class
                );
        Assertions.assertThat(getMediaSecondTimeResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(getMediaSecondTimeResponse.getBody()).isNotNull();
    }
    
    @Test
    public void SendRecommendation_MediaShouldBeAddedToLibrary() {
        /* Create user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Create Authorization header with JWT token */
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        
        /* Data */
        var user = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media in database
        int nbExistingMediaInUserLibrary = RandomUtils.nextInt(3, 7);
        for (int i = 0; i < nbExistingMediaInUserLibrary; i++) { // add some media to user library
            dummyDataCreator.addMediaToLibrary(user, dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE));
        }
        List<User> receivers = new ArrayList<>();
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) { // add some media to user library
            receivers.add(dummyDataCreator.generateUser(true));
        }
        
        /* Send recommendation */
        var sendRecommendationRequest = new RecommendationCreateRequest(
                media.getId(),
                receivers.stream().map(User::getId).toList(),
                null,
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5)
        );
        ResponseEntity<RecommendationDto> sendRecommendationResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(sendRecommendationRequest, headers),
                        RecommendationDto.class
                );
        Assertions.assertThat(sendRecommendationResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library and check that media is included */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse.getBody()).isNotNull();
        List<MediaShortDto> library = listLibraryResponse.getBody().stream().toList();
        Assertions.assertThat(library).hasSize(nbExistingMediaInUserLibrary + 1);
        var mediaFromList = library.stream()
                .filter(mediaDto -> Objects.equals(mediaDto.getId(), media.getId().toString()))
                .findAny();
        Assertions.assertThat(mediaFromList.isPresent()).isTrue();
    }

    @Test
    public void SendRecommendationForMediaAlreadyInLibrary_MediaShouldKeepPersonalFields() {
        /* Create user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Create Authorization header with JWT token */
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());

        /* Data */
        var user = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media in database
        int nbExistingMediaInUserLibrary = RandomUtils.nextInt(3, 7);
        for (int i = 0; i < nbExistingMediaInUserLibrary; i++) { // add some media to user library
            dummyDataCreator.addMediaToLibrary(user, dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE));
        }
        dummyDataCreator.addMediaToLibrary(user, media);
        List<User> receivers = new ArrayList<>();
        for (int i = 0; i < RandomUtils.nextInt(2, 5); i++) { // add some media to user library
            receivers.add(dummyDataCreator.generateUser(true));
        }

        /* Send recommendation */
        var sendRecommendationRequest = new RecommendationCreateRequest(
                media.getId(),
                receivers.stream().map(User::getId).toList(),
                null,
                RandomStringUtils.random(20, true, false),
                RandomUtils.nextInt(1, 5)
        );
        ResponseEntity<RecommendationDto> sendRecommendationResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(sendRecommendationRequest, headers),
                        RecommendationDto.class
                );
        Assertions.assertThat(sendRecommendationResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* List medias from library and check that media is included */
        ResponseEntity<CustomPageImpl<MediaShortDto>> listLibraryResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.libraryUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(listLibraryResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(listLibraryResponse.getBody()).isNotNull();
        List<MediaShortDto> library = listLibraryResponse.getBody().stream().toList();
        Assertions.assertThat(library).hasSize(nbExistingMediaInUserLibrary + 1);
        var mediaFromList = library.stream()
                .filter(mediaDto -> Objects.equals(mediaDto.getId(), media.getId().toString()))
                .findAny();
        Assertions.assertThat(mediaFromList.isPresent()).isTrue();
    }
}
