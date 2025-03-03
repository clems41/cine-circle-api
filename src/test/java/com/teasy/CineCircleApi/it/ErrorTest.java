package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationCreateRequest;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ErrorResponse;
import com.teasy.CineCircleApi.utils.HttpUtils;
import com.teasy.CineCircleApi.utils.RandomUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;

import java.time.LocalDateTime;
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
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(recommendationCreateRequest, headers),
                        ErrorResponse.class
                );
        // check response
        var expectedErrorDetails = ErrorDetails.ERR_USER_NOT_FOUND.addingArgs(nonExistingReceiverId);
        checkResponseAndInDatabase(expectedErrorDetails, response, false);
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
        var emailsInDatabase = emailRepository.findAll()
                .stream().filter(email -> Objects.equals(email.getReceiver(), existingUser.getEmail())).toList();
        Assertions.assertThat(emailsInDatabase).hasSize(1);
        var emailStored = emailsInDatabase.getFirst();
        var expectedErrorDetails = ErrorDetails.ERR_EMAIL_SENDING_REQUEST.addingArgs(emailStored.getId(), existingUser.getEmail());
        checkResponseAndInDatabase(expectedErrorDetails, response, true);
    }

    @Test
    public void checkThatErrorResponseIsCorrectWhenValidAnnotationThrowException_LengthAnnotation() {
        /* Send request to sign up new user with wring fields in request */
        var username = RandomUtils.randomString(10).toLowerCase();
        var email = String.format("%s@%s.com",
                RandomUtils.randomString(20),
                RandomUtils.randomString(5));
        var password = RandomUtils.randomString(4); // will throw error because 4 < 6
        var displayName = RandomUtils.randomString(20);
        var request = new AuthSignUpRequest(username, email, password, displayName);
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("sign-up")),
                        HttpMethod.POST,
                        new HttpEntity<>(request, null),
                        ErrorResponse.class
                );

        // check response
        var expectedErrorDetails = ErrorDetails.ERR_USER_PASSWORD_TOO_SHORT.addingArgs(password);
        checkResponseAndInDatabase(expectedErrorDetails, response, true);
    }

    @Test
    public void checkThatErrorResponseIsCorrectWhenValidAnnotationThrowException_EmailAnnotation() {
        /* Send request to sign up new user with wring fields in request */
        var username = RandomUtils.randomString(10).toLowerCase();
        var email = RandomUtils.randomString(10); // email will trigger error because it's not email format
        var password = RandomUtils.randomString(15);
        var displayName = RandomUtils.randomString(20);
        var request = new AuthSignUpRequest(username, email, password, displayName);
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("sign-up")),
                        HttpMethod.POST,
                        new HttpEntity<>(request, null),
                        ErrorResponse.class
                );

        // check response
        var expectedErrorDetails = ErrorDetails.ERR_USER_EMAIL_INCORRECT.addingArgs(email);
        checkResponseAndInDatabase(expectedErrorDetails, response, true);
    }

    @Test
    public void checkThatErrorResponseIsCorrectWhenValidAnnotationThrowException_Min() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var receiver = dummyDataCreator.generateUser(true);
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var comment = RandomUtils.randomString(50);
        var rating = 0; // will trigger an error because 0 < 1

        /* Send recommendation */
        var userIds = List.of(receiver.getId());

        var recommendationCreateRequest = new RecommendationCreateRequest(
                media.getId(), userIds, List.of(), comment, rating);
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(recommendationCreateRequest, headers),
                        ErrorResponse.class
                );

        // check response
        var expectedErrorDetails = ErrorDetails.ERR_RECOMMENDATION_RATING_INCORRECT.addingArgs(rating);
        checkResponseAndInDatabase(expectedErrorDetails, response, true);
    }

    @Test
    public void checkThatErrorResponseIsCorrectWhenValidAnnotationThrowException_Size() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var media = dummyDataCreator.generateMedia(true, MediaTypeEnum.MOVIE);
        var comment = RandomUtils.randomString(50);
        var rating = 3;

        /* Send recommendation */
        List<UUID> userIds = List.of(); // will trigger an error because empty
        var recommendationCreateRequest = new RecommendationCreateRequest(
                media.getId(), userIds, List.of(), comment, rating);
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.recommendationUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(recommendationCreateRequest, headers),
                        ErrorResponse.class
                );

        // check response
        var expectedErrorDetails = ErrorDetails.ERR_RECOMMENDATION_USER_IDS_INCORRECT;
        checkResponseAndInDatabase(expectedErrorDetails, response, true);
    }
    @Test
    public void getCircleWithBadUuid() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());

        /* Get circle with bad format Id */
        var wrongUuid = "toto";
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(wrongUuid),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        ErrorResponse.class
                );
        var expectedErrorDetails = ErrorDetails.ERR_GLOBAL_INVALID_PARAMETER.addingArgs("circle_id", wrongUuid);
        checkResponseAndInDatabase(expectedErrorDetails, response, true);

    }

    private void checkResponseAndInDatabase(ErrorDetails expectedErrorDetails,
                                            ResponseEntity<ErrorResponse> response,
                                            boolean checkStackTrace) {
        Assertions.assertThat(response.getStatusCode()).isEqualTo(expectedErrorDetails.getHttpStatus());
        Assertions.assertThat(response.getBody()).isNotNull();
        Assertions.assertThat(response.getBody().errorMessage())
                .isEqualTo(expectedErrorDetails.getMessage());
        Assertions.assertThat(response.getBody().errorCode())
                .isEqualTo(expectedErrorDetails.getCode());
        Assertions.assertThat(response.getBody().errorOnObject())
                .isEqualTo(expectedErrorDetails.getErrorOnObject().name());
        if (expectedErrorDetails.getErrorOnField() != null) {
            Assertions.assertThat(response.getBody().errorOnField())
                    .isEqualTo(expectedErrorDetails.getErrorOnField().name());
        } else {
            Assertions.assertThat(response.getBody().errorOnField()).isNull();
        }
        if (checkStackTrace) {
            Assertions.assertThat(response.getBody().errorCause()).isNotNull();
            Assertions.assertThat(response.getBody().errorStack()).isNotNull();
        }

        /* Check data in database that all recommendations have been created */
        var errors = errorRepository.findAll().stream()
                .filter(error -> Objects.equals(error.getMessage(), expectedErrorDetails.getMessage())).toList();
        Assertions.assertThat(errors).isNotEmpty();
        var error = errors.stream().toList().getFirst();
        Assertions.assertThat(error.getHttpStatusCode())
                .isEqualTo(expectedErrorDetails.getHttpStatus().value());
        Assertions.assertThat(error.getMessage())
                .isEqualTo(expectedErrorDetails.getMessage());
        Assertions.assertThat(error.getCode())
                .isEqualTo(expectedErrorDetails.getCode());
        Assertions.assertThat(error.getObject())
                .isEqualTo(expectedErrorDetails.getErrorOnObject().name());
        if (expectedErrorDetails.getErrorOnField() != null) {
            Assertions.assertThat(error.getField())
                    .isEqualTo(expectedErrorDetails.getErrorOnField().name());
        } else {
            Assertions.assertThat(error.getField()).isNull();
        }
        if (checkStackTrace) {
            Assertions.assertThat(error.getCause()).isNotNull();
            Assertions.assertThat(error.getFirstElementOfStackTrace()).isNotNull();
        }
        Assertions.assertThat(error.getTriggeredAt()).isBefore(LocalDateTime.now());
    }
}
