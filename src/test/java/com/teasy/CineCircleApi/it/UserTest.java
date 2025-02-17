package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.AuthRefreshTokenRequest;
import com.teasy.CineCircleApi.models.dtos.responses.AuthRefreshTokenResponse;
import com.teasy.CineCircleApi.models.dtos.responses.AuthSignInResponse;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ErrorResponse;
import com.teasy.CineCircleApi.utils.CustomPageImpl;
import com.teasy.CineCircleApi.utils.HttpUtils;
import com.teasy.CineCircleApi.utils.RandomUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.*;

import java.time.Instant;
import java.time.LocalDateTime;
import java.util.*;

public class UserTest extends IntegrationTestAbstract {
    @Test
    public void AddUserInRelatedUsers_Success() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var userToAdd1 = dummyDataCreator.generateUser(true);
        var userToAdd2 = dummyDataCreator.generateUser(true);
        var userToAdd3 = dummyDataCreator.generateUser(true);

        /* Add user1 as relatedUsers */
        String urlToAddUser1 = String.format("%s%srelated/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.userUrl,
                userToAdd1.getId().toString()
        );
        ResponseEntity<UserFullInfoDto> response = this.restTemplate
                .exchange(
                        urlToAddUser1,
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );
        Assertions.assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(response.getBody()).isNotNull();

        // check that user1 is in relatedUsers
        var relatedUsers = getRelatedUsers(headers, null);
        Assertions.assertThat(relatedUsers.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())))
                .isTrue();

        /* Add user2 as relatedUsers */
        String urlToAddUser2 = String.format("%s%srelated/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.userUrl,
                userToAdd2.getId().toString()
        );
        ResponseEntity<UserFullInfoDto> response2 = this.restTemplate
                .exchange(
                        urlToAddUser2,
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );
        Assertions.assertThat(response2.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(response2.getBody()).isNotNull();

        // check that user1 is in relatedUsers
        var relatedUsers2 = getRelatedUsers(headers, null);
        Assertions.assertThat(relatedUsers2).hasSize(2);
        Assertions.assertThat(relatedUsers2.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd2.getId().toString())))
                .isTrue();

        /* Add user1 as relatedUsers */
        String urlToAddUser3 = String.format("%s%srelated/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.userUrl,
                userToAdd3.getId().toString()
        );
        ResponseEntity<UserFullInfoDto> response3 = this.restTemplate
                .exchange(
                        urlToAddUser3,
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );
        Assertions.assertThat(response3.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(response3.getBody()).isNotNull();

        // check that user1 is in relatedUsers
        var relatedUsers3 = getRelatedUsers(headers, null);
        Assertions.assertThat(relatedUsers3).hasSize(3);
        Assertions.assertThat(relatedUsers3.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd3.getId().toString())))
                .isTrue();

        /* Check that all 3 users are in relatedUsers */
        Assertions.assertThat(relatedUsers3.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())))
                .isTrue();
        Assertions.assertThat(relatedUsers3.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd2.getId().toString())))
                .isTrue();
        Assertions.assertThat(relatedUsers3.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd3.getId().toString())))
                .isTrue();
    }

    @Test
    public void AddUserInRelatedUsers_NotExistingUser() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var fakeUser = dummyDataCreator.generateUser(false);
        fakeUser.setId(UUID.randomUUID());

        /* Add fakeUser as relatedUsers */
        String urlToAddFakeUser = String.format("%s%srelated/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.userUrl,
                fakeUser.getId().toString()
        );
        ResponseEntity<UserFullInfoDto> response = this.restTemplate
                .exchange(
                        urlToAddFakeUser,
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );
        Assertions.assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
    }

    @Test
    public void RemoveUserFromRelatedUsers_Success() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var userToRemove = dummyDataCreator.generateUser(true);
        var userToAdd1 = dummyDataCreator.generateUser(true);
        var userToAdd2 = dummyDataCreator.generateUser(true);
        authenticatedUser.addRelatedUser(userToRemove);
        authenticatedUser.addRelatedUser(userToAdd1);
        authenticatedUser.addRelatedUser(userToAdd2);
        userRepository.save(authenticatedUser);

        /* Remove userToRemove from relatedUsers */
        String urlToRemoveUser = String.format("%s%srelated/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.userUrl,
                userToRemove.getId().toString()
        );
        ResponseEntity<UserFullInfoDto> response = this.restTemplate
                .exchange(
                        urlToRemoveUser,
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );
        Assertions.assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(response.getBody()).isNotNull();

        // check that userToRemove is not in relatedUsers
        var relatedUsers = getRelatedUsers(headers, null);
        Assertions.assertThat(relatedUsers).hasSize(2);
        Assertions.assertThat(relatedUsers.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToRemove.getId().toString())))
                .isFalse();

        /* Check that users 1 and 2 are in relatedUsers */
        Assertions.assertThat(relatedUsers.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())))
                .isTrue();
        Assertions.assertThat(relatedUsers.stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd2.getId().toString())))
                .isTrue();
    }

    @Test
    public void RemoveUserFromRelatedUsers_NotExistingUser() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var fakeUser = dummyDataCreator.generateUser(false);
        fakeUser.setId(UUID.randomUUID());

        /* Add fakeUser as relatedUsers */
        String urlToAddFakeUser = String.format("%s%srelated/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.userUrl,
                fakeUser.getId().toString()
        );
        ResponseEntity<UserFullInfoDto> response = this.restTemplate
                .exchange(
                        urlToAddFakeUser,
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );
        Assertions.assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
    }

    @Test
    public void refreshToken() {
        /* Sign up user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Sign in user ang get token */
        ResponseEntity<AuthSignInResponse> signInResponse = this.restTemplate
                .withBasicAuth(signUpRequest.username(), signUpRequest.password())
                .getForEntity(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("sign-in")),
                        AuthSignInResponse.class
                );
        Assertions.assertThat(signInResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(signInResponse.getBody()).isNotNull();
        var jwtToken = signInResponse.getBody().getToken().tokenString();
        Assertions.assertThat(jwtToken).isNotEmpty();
        var jwtRefreshToken = signInResponse.getBody().getRefreshToken();
        Assertions.assertThat(jwtRefreshToken.tokenString()).isNotEmpty();
        Assertions.assertThat(jwtRefreshToken.expirationDate()).isAfter(LocalDateTime.now());

        /* Refresh token */
        var refreshRequest = new AuthRefreshTokenRequest(jwtToken, jwtRefreshToken.tokenString());
        ResponseEntity<AuthRefreshTokenResponse> refreshTokenResponse = this.restTemplate
                .postForEntity(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("refresh")),
                        refreshRequest,
                        AuthRefreshTokenResponse.class
                );
        Assertions.assertThat(refreshTokenResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(refreshTokenResponse.getBody()).isNotNull();
        var newJwtToken = refreshTokenResponse.getBody().getJwtToken();
        Assertions.assertThat(newJwtToken.tokenString()).isNotEmpty();
        Assertions.assertThat(newJwtToken.expirationDate()).isAfter(Instant.now());
    }

    @Test
    public void refreshToken_badRefreshToken() {
        /* Sign up user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Sign in user ang get token */
        ResponseEntity<AuthSignInResponse> signInResponse = this.restTemplate
                .withBasicAuth(signUpRequest.username(), signUpRequest.password())
                .getForEntity(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("sign-in")),
                        AuthSignInResponse.class
                );
        Assertions.assertThat(signInResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(signInResponse.getBody()).isNotNull();
        var jwtToken = signInResponse.getBody().getToken().tokenString();
        Assertions.assertThat(jwtToken).isNotEmpty();
        var jwtRefreshToken = signInResponse.getBody().getRefreshToken();
        Assertions.assertThat(jwtRefreshToken.tokenString()).isNotEmpty();
        Assertions.assertThat(jwtRefreshToken.expirationDate()).isAfter(LocalDateTime.now());

        /* Refresh token */
        var refreshRequest = new AuthRefreshTokenRequest(jwtToken, "toto");
        ResponseEntity<ErrorResponse> refreshTokenResponse = this.restTemplate
                .postForEntity(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("refresh")),
                        refreshRequest,
                        ErrorResponse.class
                );
        Assertions.assertThat(refreshTokenResponse.getStatusCode()).isEqualTo(HttpStatus.UNAUTHORIZED);
        Assertions.assertThat(refreshTokenResponse.getBody()).isNotNull();
        Assertions.assertThat(refreshTokenResponse.getBody().errorCode()).isEqualTo(ErrorDetails.ERR_AUTH_CANNOT_REFRESH_TOKEN.getCode());
    }

    @Test
    public void refreshToken_badToken() {
        /* Sign up user */
        var signUpRequest = authenticator.authenticateNewUser();

        /* Sign in user ang get token */
        ResponseEntity<AuthSignInResponse> signInResponse = this.restTemplate
                .withBasicAuth(signUpRequest.username(), signUpRequest.password())
                .getForEntity(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("sign-in")),
                        AuthSignInResponse.class
                );
        Assertions.assertThat(signInResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(signInResponse.getBody()).isNotNull();
        var jwtToken = signInResponse.getBody().getToken().tokenString();
        Assertions.assertThat(jwtToken).isNotEmpty();
        var jwtRefreshToken = signInResponse.getBody().getRefreshToken();
        Assertions.assertThat(jwtRefreshToken.tokenString()).isNotEmpty();
        Assertions.assertThat(jwtRefreshToken.expirationDate()).isAfter(LocalDateTime.now());

        /* Refresh token with bad token */
        var refreshRequest = new AuthRefreshTokenRequest("toto", jwtRefreshToken.tokenString());
        ResponseEntity<ErrorResponse> refreshTokenResponse = this.restTemplate
                .postForEntity(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.authUrl.concat("refresh")),
                        refreshRequest,
                        ErrorResponse.class
                );
        Assertions.assertThat(refreshTokenResponse.getStatusCode()).isEqualTo(HttpStatus.UNAUTHORIZED);
        Assertions.assertThat(refreshTokenResponse.getBody()).isNotNull();
        Assertions.assertThat(refreshTokenResponse.getBody().errorCode()).isEqualTo(ErrorDetails.ERR_AUTH_CANNOT_REFRESH_TOKEN.getCode());
    }

    @Test
    public void defaultSortingForRelatedUsers() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();

        // create relatedUsers with specific number of received recommendations from authenticated user
        var userWhoReceived5Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived4Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived3Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived2Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived0Recommendations = dummyDataCreator.generateUser(true);
        authenticatedUser.addRelatedUser(userWhoReceived5Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived4Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived3Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived2Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived0Recommendations);
        userRepository.save(authenticatedUser);

        // create recommendations for each related user
        for (int i = 0; i < 5; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived5Recommendations, null, true);
        }
        for (int i = 0; i < 4; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived4Recommendations, null, true);
        }
        for (int i = 0; i < 3; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived3Recommendations, null, true);
        }
        for (int i = 0; i < 2; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived2Recommendations, null, true);
        }

        // create recommendations between other users
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived4Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived5Recommendations, userWhoReceived4Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived2Recommendations, userWhoReceived3Recommendations, null, true);


        /* Get related users */
        Map<String, Object> queryParameters = new HashMap<>();
        queryParameters.put("page", 0);
        queryParameters.put("size", 10);
        var relatedUsers = getRelatedUsers(headers, queryParameters);

        /* Check response and sorting */
        Assertions.assertThat(relatedUsers).isNotNull();
        Assertions.assertThat(relatedUsers).hasSize(5);
        Assertions.assertThat(relatedUsers.getFirst().getId()).isEqualTo(userWhoReceived5Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(1).getId()).isEqualTo(userWhoReceived4Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(2).getId()).isEqualTo(userWhoReceived3Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(3).getId()).isEqualTo(userWhoReceived2Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.getLast().getId()).isEqualTo(userWhoReceived0Recommendations.getId().toString());
    }

    @Test
    public void specificSortingForRelatedUsers() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();

        // create relatedUsers with specific number of received recommendations from authenticated user
        var relatedUser1 = dummyDataCreator.generateUser(true);
        var relatedUser2 = dummyDataCreator.generateUser(true);
        var relatedUser3 = dummyDataCreator.generateUser(true);
        var relatedUser4 = dummyDataCreator.generateUser(true);
        var relatedUser5 = dummyDataCreator.generateUser(true);
        authenticatedUser.addRelatedUser(relatedUser1);
        authenticatedUser.addRelatedUser(relatedUser2);
        authenticatedUser.addRelatedUser(relatedUser3);
        authenticatedUser.addRelatedUser(relatedUser4);
        authenticatedUser.addRelatedUser(relatedUser5);
        userRepository.save(authenticatedUser);

        /* Get related users */
        Map<String, Object> queryParameters = new HashMap<>();
        queryParameters.put("page", 0);
        queryParameters.put("size", 10);
        queryParameters.put("sort", "username,desc");
        var relatedUsers = getRelatedUsers(headers, queryParameters);

        /* Check response and sorting by username desc */
        Assertions.assertThat(relatedUsers).isNotNull();
        Assertions.assertThat(relatedUsers).hasSize(5);
        var expectedRelatedUsers = new ArrayList<>(List
                .of(relatedUser1, relatedUser2, relatedUser3, relatedUser4, relatedUser5));
        expectedRelatedUsers.sort(Comparator.comparing(User::getUsername).reversed());
        for (int i = 0; i < 5; i++) {
            Assertions.assertThat(relatedUsers.get(i).getId()).isEqualTo(expectedRelatedUsers.get(i).getId().toString());
        }
    }

    @Test
    public void emptyQueryAndDefaultSortingForRelatedUsers() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();

        // create relatedUsers with specific number of received recommendations from authenticated user
        var userWhoReceived5Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived4Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived3Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived2Recommendations = dummyDataCreator.generateUser(true);
        var userWhoReceived0Recommendations = dummyDataCreator.generateUser(true);
        authenticatedUser.addRelatedUser(userWhoReceived5Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived4Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived3Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived2Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived0Recommendations);
        userRepository.save(authenticatedUser);

        // create recommendations for each related user
        for (int i = 0; i < 5; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived5Recommendations, null, true);
        }
        for (int i = 0; i < 4; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived4Recommendations, null, true);
        }
        for (int i = 0; i < 3; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived3Recommendations, null, true);
        }
        for (int i = 0; i < 2; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived2Recommendations, null, true);
        }

        // create recommendations between other users
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived4Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived3Recommendations, userWhoReceived2Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived5Recommendations, userWhoReceived4Recommendations, null, true);
        dummyDataCreator.generateRecommendation(true, userWhoReceived2Recommendations, userWhoReceived3Recommendations, null, true);


        /* Get related users */
        Map<String, Object> queryParameters = new HashMap<>();
        queryParameters.put("page", 0);
        queryParameters.put("size", 10);
        queryParameters.put("query", "");
        var relatedUsers = getRelatedUsers(headers, queryParameters);

        /* Check response and sorting */
        Assertions.assertThat(relatedUsers).isNotNull();
        Assertions.assertThat(relatedUsers).hasSize(5);
        Assertions.assertThat(relatedUsers.getFirst().getId()).isEqualTo(userWhoReceived5Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(1).getId()).isEqualTo(userWhoReceived4Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(2).getId()).isEqualTo(userWhoReceived3Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(3).getId()).isEqualTo(userWhoReceived2Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.getLast().getId()).isEqualTo(userWhoReceived0Recommendations.getId().toString());
    }

    @Test
    public void queryAndDefaultSortingForRelatedUsers() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();

        // create word that will use for matching username
        var query = RandomUtils.randomString(4);

        // create relatedUsers with specific number of received recommendations from authenticated user
        var userWhoReceived5Recommendations = dummyDataCreator.generateUserWithUsername(true, RandomUtils.randomString(4).concat(query).concat(RandomUtils.randomString(2)));
        var userWhoReceived4Recommendations = dummyDataCreator.generateUserWithUsername(true, RandomUtils.randomString(4).concat(query).concat(RandomUtils.randomString(2)));
        var userWhoReceived3Recommendations = dummyDataCreator.generateUserWithUsername(true, RandomUtils.randomString(4).concat(query).concat(RandomUtils.randomString(2)));
        var userWhoReceived2Recommendations = dummyDataCreator.generateUserWithUsername(true, RandomUtils.randomString(4).concat(query).concat(RandomUtils.randomString(2)));
        var userWhoReceived1Recommendations = dummyDataCreator.generateUserWithUsername(true, RandomUtils.randomString(4).concat(query).concat(RandomUtils.randomString(2)));
        authenticatedUser.addRelatedUser(userWhoReceived5Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived4Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived3Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived2Recommendations);
        authenticatedUser.addRelatedUser(userWhoReceived1Recommendations);

        // create non-matching relatedUser
        for (int i = 0; i < 10; i++) {
            var user = dummyDataCreator.generateUser(true);
            authenticatedUser.addRelatedUser(user);
        }

        userRepository.save(authenticatedUser);

        // create recommendations for each related user
        for (int i = 0; i < 5; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived5Recommendations, null, true);
        }
        for (int i = 0; i < 4; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived4Recommendations, null, true);
        }
        for (int i = 0; i < 3; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived3Recommendations, null, true);
        }
        for (int i = 0; i < 2; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived2Recommendations, null, true);
        }
        for (int i = 0; i < 1; i++) {
            dummyDataCreator.generateRecommendation(true, authenticatedUser, userWhoReceived1Recommendations, null, true);
        }


        /* Get related users */
        Map<String, Object> queryParameters = new HashMap<>();
        queryParameters.put("page", 0);
        queryParameters.put("size", 10);
        queryParameters.put("query", query);
        var relatedUsers = getRelatedUsers(headers, queryParameters);

        /* Check response and sorting */
        Assertions.assertThat(relatedUsers).isNotNull();
        Assertions.assertThat(relatedUsers).hasSize(5);
        Assertions.assertThat(relatedUsers.getFirst().getId()).isEqualTo(userWhoReceived5Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(1).getId()).isEqualTo(userWhoReceived4Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(2).getId()).isEqualTo(userWhoReceived3Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.get(3).getId()).isEqualTo(userWhoReceived2Recommendations.getId().toString());
        Assertions.assertThat(relatedUsers.getLast().getId()).isEqualTo(userWhoReceived1Recommendations.getId().toString());
    }
    
    private List<UserDto> getRelatedUsers(HttpHeaders headers, Map<String, Object> queryParameters) {
        ResponseEntity<CustomPageImpl<UserDto>> response = this.restTemplate
                .exchange(
                        HttpUtils.getUriWithQueryParameter(port, HttpUtils.userUrl.concat("related"), queryParameters),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>() {
                        }
                );
        Assertions.assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(response.getBody()).isNotNull();
        Assertions.assertThat(response.getBody().getContent()).isNotNull();
        return response.getBody().getContent();
    }
}
