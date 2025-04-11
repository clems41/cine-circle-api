package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.MediaFullDto;
import com.teasy.CineCircleApi.models.dtos.MediaShortDto;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.utils.CustomPageImpl;
import com.teasy.CineCircleApi.utils.HttpUtils;
import com.teasy.CineCircleApi.utils.RandomUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import java.util.List;
import java.util.Objects;
import java.util.UUID;

public class WatchlistTest extends IntegrationTestAbstract {
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
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(media1.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMedia1Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Add non-existing media to watchlist */
        ResponseEntity<String> addNonExistingMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(nonExistingMediaId.toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Add media2 to watchlist */
        ResponseEntity<String> addMedia2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(media2.getId().toString()),
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
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(nonExistingMediaId.toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeNonExistingMediaResponse.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);

        /* Remove media2 from watchlist */
        ResponseEntity<String> removeMedia2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(media2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeMedia2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Remove media2 a 2nd time from watchlist */
        ResponseEntity<String> removeMedia2SecondTimeResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(media2.getId().toString()),
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
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(media3.getId().toString()),
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

    @Test
    public void AddMedia_CheckThatIsInWatchlistIsTrueWhenGettingMedia() {
        /* Init */
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE); // create media in database

        /* Create user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Create Authorization header with JWT token */
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());

        /* Add media to watchlist */
        ResponseEntity<String> addMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.watchlistUrl).concat("/").concat(media.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addMediaResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get media */
        ResponseEntity<MediaFullDto> getMediaResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.mediaUrl).concat(media.getId().toString()),
                        HttpMethod.GET,
                        new HttpEntity<>(headers),
                        MediaFullDto.class
                );
        Assertions.assertThat(getMediaResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(getMediaResponse.getBody()).isNotNull();
        Assertions.assertThat(getMediaResponse.getBody().getIsInWatchlist()).isTrue();

        /* Add media to relatedUser's heading */
        var relatedUserUsername = RandomUtils.randomString(10);
        var relatedUser = dummyDataCreator.generateUserWithUsername(true, relatedUserUsername);
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        authenticatedUser.addRelatedUser(relatedUser);
        userRepository.save(authenticatedUser);
        relatedUser.addMediaToHeadings(media);
        userRepository.save(relatedUser);

        /* Get headings for related user and check boolean */
        ResponseEntity<List<MediaShortDto>> headingsResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.headingsUrl).concat("users/").concat(relatedUser.getId().toString()),
                        HttpMethod.GET,
                        new HttpEntity<>(headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(headingsResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(headingsResponse.getBody()).isNotNull();
        List<MediaShortDto> headings = headingsResponse.getBody();
        Assertions.assertThat(headings).hasSize(1);
        Assertions.assertThat(headings.getFirst().getIsInWatchlist()).isTrue();

        /* RelatedUser send recommendation to authenticatedUser for this media */
        var recommendation = new Recommendation(UUID.randomUUID(), relatedUser, media, authenticatedUser, "top", 5);
        recommendationRepository.save(recommendation);

        /* Get recommendations and check boolean */
        ResponseEntity<CustomPageImpl<RecommendationDto>> recommendationsResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.GET,
                        new HttpEntity<>(headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(recommendationsResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(recommendationsResponse.getBody()).isNotNull();
        Assertions.assertThat(recommendationsResponse.getBody().getContent()).hasSize(1);
        Assertions.assertThat(recommendationsResponse.getBody().getContent().getFirst().getMedia().getIsInWatchlist()).isTrue();
    }
}
