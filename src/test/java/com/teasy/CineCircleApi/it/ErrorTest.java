package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ErrorResponse;
import com.teasy.CineCircleApi.utils.HttpUtils;
import com.teasy.CineCircleApi.utils.RandomUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.http.*;

import java.util.*;

public class ErrorTest extends IntegrationTestAbstract {
    @Test
    public void checkErrorIsStoredInDatabase() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var nonExistingReceiverId = UUID.randomUUID();
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var comment = RandomUtils.randomString(50);
        var rating = RandomUtils.randomInt(1, 5);

        /* Send recommendation with a non-existing user as receiver, should throw error */
        var userIds = List.of(nonExistingReceiverId);
        var recommendationCreateRequest = new RecommendationCreateRequest(
                media.getId(), userIds, List.of(), comment, rating);
        ResponseEntity<ErrorResponse> recommendationCreateResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(recommendationCreateRequest, headers),
                        ErrorResponse.class
                );
        // check response
        var expectedErrorDetails = ErrorDetails.ERR_USER_NOT_FOUND.addingArgs(nonExistingReceiverId);
        Assertions.assertThat(recommendationCreateResponse.getStatusCode()).isEqualTo(expectedErrorDetails.getHttpStatus());
        Assertions.assertThat(recommendationCreateResponse.getBody()).isNotNull();
        Assertions.assertThat(recommendationCreateResponse.getBody().errorMessage())
                .isEqualTo(expectedErrorDetails.getMessage());
        Assertions.assertThat(recommendationCreateResponse.getBody().errorCode())
                .isEqualTo(expectedErrorDetails.getCode());
        Assertions.assertThat(recommendationCreateResponse.getBody().errorOnObject())
                .isEqualTo(expectedErrorDetails.getErrorOnObject().name());
        Assertions.assertThat(recommendationCreateResponse.getBody().errorOnField())
                .isEqualTo(expectedErrorDetails.getErrorOnField().name());
        Assertions.assertThat(recommendationCreateResponse.getBody().errorCause()).isNull();
        Assertions.assertThat(recommendationCreateResponse.getBody().errorStack()).isNull();

        /* Check data in database that all recommendations have been created */
        var errors = errorRepository.findAll().stream()
                .filter(error -> Objects.equals(error.getCode(), expectedErrorDetails.getCode())).toList();
        Assertions.assertThat(errors).hasSize(1);
        var error = errors.stream().toList().getFirst();
        Assertions.assertThat(error.getHttpStatusCode())
                .isEqualTo(expectedErrorDetails.getHttpStatus().value());
        Assertions.assertThat(error.getMessage())
                .isEqualTo(expectedErrorDetails.getMessage());
        Assertions.assertThat(error.getCode())
                .isEqualTo(expectedErrorDetails.getCode());
        Assertions.assertThat(error.getObject())
                .isEqualTo(expectedErrorDetails.getErrorOnObject().name());
        Assertions.assertThat(error.getField())
                .isEqualTo(expectedErrorDetails.getErrorOnField().name());
        Assertions.assertThat(error.getCause()).isNull();
        Assertions.assertThat(error.getFirstElementOfStackTrace()).isNull();
    }

    @Test
    public void checkErrorIsStoredInDatabase_withStackTrace() {
        /* Data */
        var existingUser = dummyDataCreator.generateUser(true);

        /* Send request to reset password via email, but as the SMTP is not configured for the tests, it should throw an error with cause */
        Map<String, Object> queryParams = new HashMap<>();
        queryParams.put("email", existingUser.getEmail());
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getUriWithQueryParameter(port, HttpUtils.userUrl.concat("reset-password"), queryParams),
                        HttpMethod.GET,
                        new HttpEntity<>(null, null),
                        ErrorResponse.class
                );
        // check response
        var expectedErrorDetails = ErrorDetails.ERR_EMAIL_SENDING_REQUEST.addingArgs(existingUser.getEmail());
        Assertions.assertThat(response.getStatusCode()).isEqualTo(expectedErrorDetails.getHttpStatus());
        Assertions.assertThat(response.getBody()).isNotNull();
        Assertions.assertThat(response.getBody().errorMessage())
                .isEqualTo(expectedErrorDetails.getMessage());
        Assertions.assertThat(response.getBody().errorCode())
                .isEqualTo(expectedErrorDetails.getCode());
        Assertions.assertThat(response.getBody().errorOnObject())
                .isEqualTo(expectedErrorDetails.getErrorOnObject().name());
        Assertions.assertThat(response.getBody().errorOnField()).isNull(); // this field is null for this error
        Assertions.assertThat(response.getBody().errorCause()).isNotNull();
        Assertions.assertThat(response.getBody().errorStack()).isNotNull();

        /* Check data in database that all recommendations have been created */
        var errors = errorRepository.findAll().stream()
                .filter(error -> Objects.equals(error.getCode(), expectedErrorDetails.getCode())).toList();
        Assertions.assertThat(errors).hasSize(1);
        var error = errors.stream().toList().getFirst();
        Assertions.assertThat(error.getHttpStatusCode())
                .isEqualTo(expectedErrorDetails.getHttpStatus().value());
        Assertions.assertThat(error.getMessage())
                .isEqualTo(expectedErrorDetails.getMessage());
        Assertions.assertThat(error.getCode())
                .isEqualTo(expectedErrorDetails.getCode());
        Assertions.assertThat(error.getObject())
                .isEqualTo(expectedErrorDetails.getErrorOnObject().name());
        Assertions.assertThat(error.getField()).isNull(); // this field is null for this error
        Assertions.assertThat(error.getCause()).isNotNull();
        Assertions.assertThat(error.getFirstElementOfStackTrace()).isNotNull();
    }
}
