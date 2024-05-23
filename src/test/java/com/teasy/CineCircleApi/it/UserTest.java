package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.utils.HttpUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import java.util.Objects;
import java.util.UUID;

public class UserTest extends IntegrationTestAbstract {
    @Test
    public void AddUserInRelatedUsers_Success() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var userToAdd1 = dummyDataCreator.generateUser(true);
        var userToAdd2 = dummyDataCreator.generateUser(true);
        var userToAdd3 = dummyDataCreator.generateUser(true);

        /* Add user1 as relatedUsers */
        String urlToAddUser1 = String.format("%s%sme/related/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl,
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
        Assertions.assertThat(response.getBody().getRelatedUsers().size()).isEqualTo(1);
        Assertions.assertThat(response.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())))
                .isTrue();

        /* Add user2 as relatedUsers */
        String urlToAddUser2 = String.format("%s%sme/related/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl,
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
        Assertions.assertThat(response2.getBody().getRelatedUsers().size()).isEqualTo(2);
        Assertions.assertThat(response2.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd2.getId().toString())))
                .isTrue();

        /* Add user1 as relatedUsers */
        String urlToAddUser3 = String.format("%s%sme/related/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl,
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
        Assertions.assertThat(response3.getBody().getRelatedUsers().size()).isEqualTo(3);
        Assertions.assertThat(response3.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd3.getId().toString())))
                .isTrue();

        /* Get user full info */
        String urlGetUserFullInfo = String.format("%s%sme",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl
        );
        ResponseEntity<UserFullInfoDto> responseUserFullInfo = this.restTemplate
                .exchange(
                        urlGetUserFullInfo,
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );

        /* Check that all 3 users are in relatedUsers */
        Assertions.assertThat(responseUserFullInfo.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(responseUserFullInfo.getBody()).isNotNull();
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().size()).isEqualTo(3);
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())))
                .isTrue();
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd2.getId().toString())))
                .isTrue();
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd3.getId().toString())))
                .isTrue();
    }

    @Test
    public void AddUserInRelatedUsers_NotExistingUser() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var fakeUser = dummyDataCreator.generateUser(false);
        fakeUser.setId(UUID.randomUUID());

        /* Add fakeUser as relatedUsers */
        String urlToAddFakeUser = String.format("%s%sme/related/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl,
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
        String urlToRemoveUser = String.format("%s%sme/related/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl,
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
        Assertions.assertThat(response.getBody().getRelatedUsers().size()).isEqualTo(2);
        Assertions.assertThat(response.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToRemove.getId().toString())))
                .isFalse();

        /* Get user full info */
        String urlGetUserFullInfo = String.format("%s%sme",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl
        );
        ResponseEntity<UserFullInfoDto> responseUserFullInfo = this.restTemplate
                .exchange(
                        urlGetUserFullInfo,
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        UserFullInfoDto.class
                );

        /* Check that only users 1 and 2 are in relatedUsers */
        Assertions.assertThat(responseUserFullInfo.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(responseUserFullInfo.getBody()).isNotNull();
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().size()).isEqualTo(2);
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())))
                .isTrue();
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToAdd2.getId().toString())))
                .isTrue();
        Assertions.assertThat(responseUserFullInfo.getBody().getRelatedUsers().stream()
                        .anyMatch(userDto -> Objects.equals(userDto.getId(), userToRemove.getId().toString())))
                .isFalse();
    }

    @Test
    public void RemoveUserFromRelatedUsers_NotExistingUser() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var fakeUser = dummyDataCreator.generateUser(false);
        fakeUser.setId(UUID.randomUUID());

        /* Add fakeUser as relatedUsers */
        String urlToAddFakeUser = String.format("%s%sme/related/%s",
                HttpUtils.getTestingUrl(port),
                HttpUtils.authUrl,
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
}
