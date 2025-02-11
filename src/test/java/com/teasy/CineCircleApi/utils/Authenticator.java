package com.teasy.CineCircleApi.utils;

import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.models.dtos.responses.AuthSignInResponse;
import org.assertj.core.api.Assertions;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

public class Authenticator {

    private final static String authSignUpUrl = "/auth/sign-up";
    private final static String authSignInUrl = "/auth/sign-in";
    private final TestRestTemplate restTemplate;
    private final int port;

    public Authenticator(
            TestRestTemplate restTemplate,
            int port
    ) {
        this.restTemplate = restTemplate;
        this.port = port;
    }

    public AuthSignUpRequest authenticateNewUser() {
        var authSignUpRequest = generateAuthSignUpRequest();
        ResponseEntity<UserFullInfoDto> signUpResponse = this.restTemplate
                .postForEntity(
                        HttpUtils.getTestingUrl(port).concat(authSignUpUrl),
                        authSignUpRequest,
                        UserFullInfoDto.class
                );
        Assertions.assertThat(signUpResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(signUpResponse.getBody()).isNotNull();
        Assertions.assertThat(signUpResponse.getBody().getUsername()).isEqualTo(authSignUpRequest.username().toLowerCase());
        Assertions.assertThat(signUpResponse.getBody().getEmail()).isEqualTo(authSignUpRequest.email());
        Assertions.assertThat(signUpResponse.getBody().getDisplayName()).isEqualTo(authSignUpRequest.displayName());
        Assertions.assertThat(signUpResponse.getBody().getTopicName()).isNotEmpty();
        Assertions.assertThat(signUpResponse.getBody().getId()).isNotEmpty();

        return authSignUpRequest;
    }

    public HttpHeaders authenticateUserAndGetHeadersWithJwtToken(String username, String password) {
        ResponseEntity<AuthSignInResponse> signInResponse = this.restTemplate
                .withBasicAuth(username, password)
                .getForEntity(
                        HttpUtils.getTestingUrl(port).concat(authSignInUrl),
                        AuthSignInResponse.class
                );
        Assertions.assertThat(signInResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(signInResponse.getBody()).isNotNull();
        var jwtToken = signInResponse.getBody().getToken().tokenString();
        Assertions.assertThat(jwtToken).isNotEmpty();

        /* Create Authorization header with JWT token */
        HttpHeaders headers = new HttpHeaders();
        headers.setBearerAuth(jwtToken);
        return headers;
    }

    private AuthSignUpRequest generateAuthSignUpRequest() {
        var username = RandomUtils.randomString(10).toLowerCase();
        var email = String.format("%s@%s.com",
                RandomUtils.randomString(20),
                RandomUtils.randomString(5));
        var password = RandomUtils.randomString(15);
        var displayName = RandomUtils.randomString(20);
        return new AuthSignUpRequest(username, email, password, displayName);
    }
}
