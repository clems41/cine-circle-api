package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.CineCircleApiApplication;
import com.teasy.CineCircleApi.models.dtos.CircleDto;
import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.repositories.*;
import com.teasy.CineCircleApi.utils.Authenticator;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import com.teasy.CineCircleApi.utils.HttpUtils;
import org.apache.commons.lang3.RandomStringUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.test.context.ActiveProfiles;

import java.util.Objects;

@ActiveProfiles("test")
@SpringBootTest(classes = CineCircleApiApplication.class, webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public class CircleTest {
    @LocalServerPort
    private int port;

    @Autowired
    private TestRestTemplate restTemplate;

    @Autowired
    private MediaRepository mediaRepository;

    @Autowired
    private CircleRepository circleRepository;

    @Autowired
    private LibraryRepository libraryRepository;

    @Autowired
    private RecommendationRepository recommendationRepository;

    @Autowired
    private UserRepository userRepository;
    private Authenticator authenticator;
    private DummyDataCreator dummyDataCreator;

    @BeforeEach
    public void setUp() {
        authenticator = new Authenticator(restTemplate, port);
        dummyDataCreator = new DummyDataCreator(userRepository, mediaRepository, recommendationRepository, libraryRepository, circleRepository);
    }

    @Test
    public void CRUDCircleAddAndRemoveUsers() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var userToAdd1 = dummyDataCreator.generateUser(true);
        var userToAdd2 = dummyDataCreator.generateUser(true);
        var userToAdd3 = dummyDataCreator.generateUser(true);
        var name1 = RandomStringUtils.random(15, true, false);
        var name2 = RandomStringUtils.random(15, true, false);
        var description1 = RandomStringUtils.random(30, true, true);
        var description2 = RandomStringUtils.random(30, true, true);

        /* Create circle and check that creator is has been added the users list */
        var createCircleRequest = new CircleCreateUpdateRequest(name1, description1, false);
        ResponseEntity<CircleDto> createCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port).concat(HttpUtils.circleUrl),
                        HttpMethod.POST,
                        new HttpEntity<>(createCircleRequest, headers),
                        CircleDto.class
                );
        Assertions.assertThat(createCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        var circle = Objects.requireNonNull(createCircleResponse.getBody());
        var circleId = circle.getId();
        Assertions.assertThat(
                circle.getUsers().stream().anyMatch(
                        userDto -> Objects.equals(userDto.getId(), authenticatedUser.getId().toString())
                )
        ).isTrue(); // check that authenticated user is in the circle
        Assertions.assertThat(circle.getName()).isEqualTo(name1);
        Assertions.assertThat(circle.getDescription()).isEqualTo(description1);
        Assertions.assertThat(circle.getIsPublic()).isFalse();

        /* Add user1 to circle */
        ResponseEntity<String> addUser1Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId)
                                .concat(HttpUtils.userUrl)
                                .concat(userToAdd1.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addUser1Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get circle and check that user1 is in the circle */
        ResponseEntity<CircleDto> getCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        CircleDto.class
                );
        Assertions.assertThat(getCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        circle = Objects.requireNonNull(getCircleResponse.getBody());
        Assertions.assertThat(
                circle.getUsers().stream().anyMatch(
                        userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())
                )
        ).isTrue(); // check that user1 is in the circle

        /* Update circle info */
        var updateCircleRequest = new CircleCreateUpdateRequest(name2, description2, true);
        ResponseEntity<CircleDto> updateCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId),
                        HttpMethod.PUT,
                        new HttpEntity<>(updateCircleRequest, headers),
                        CircleDto.class
                );
        Assertions.assertThat(updateCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get circle and check that info have been updated */
        ResponseEntity<CircleDto> getCircleResponse2 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        CircleDto.class
                );
        Assertions.assertThat(getCircleResponse2.getStatusCode()).isEqualTo(HttpStatus.OK);
        circle = Objects.requireNonNull(getCircleResponse2.getBody());
        Assertions.assertThat(circle.getName()).isEqualTo(name2);
        Assertions.assertThat(circle.getDescription()).isEqualTo(description2);
        Assertions.assertThat(circle.getIsPublic()).isTrue();

        /* Add user2 and user3 to circle */
        ResponseEntity<String> addUser2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId)
                                .concat(HttpUtils.userUrl)
                                .concat(userToAdd2.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addUser2Response.getStatusCode()).isEqualTo(HttpStatus.OK);
        ResponseEntity<String> addUser3Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId)
                                .concat(HttpUtils.userUrl)
                                .concat(userToAdd3.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addUser3Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get circle and check that user1, user2 and user3 are in the circle */
        ResponseEntity<CircleDto> getCircleResponse3 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        CircleDto.class
                );
        Assertions.assertThat(getCircleResponse3.getStatusCode()).isEqualTo(HttpStatus.OK);
        circle = Objects.requireNonNull(getCircleResponse3.getBody());
        Assertions.assertThat(
                circle.getUsers().stream().anyMatch(
                        userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())
                )
        ).isTrue(); // check that user1 is in the circle
        Assertions.assertThat(
                circle.getUsers().stream().anyMatch(
                        userDto -> Objects.equals(userDto.getId(), userToAdd2.getId().toString())
                )
        ).isTrue(); // check that user2 is in the circle
        Assertions.assertThat(
                circle.getUsers().stream().anyMatch(
                        userDto -> Objects.equals(userDto.getId(), userToAdd3.getId().toString())
                )
        ).isTrue(); // check that user3 is in the circle

        /* Remove user2 from circle */
        ResponseEntity<String> removeUser2Response = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId)
                                .concat(HttpUtils.userUrl)
                                .concat(userToAdd2.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeUser2Response.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Get circle and check that user1 and user3 are in the circle */
        ResponseEntity<CircleDto> getCircleResponse4 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        CircleDto.class
                );
        Assertions.assertThat(getCircleResponse4.getStatusCode()).isEqualTo(HttpStatus.OK);
        circle = Objects.requireNonNull(getCircleResponse4.getBody());
        Assertions.assertThat(
                circle.getUsers().stream().anyMatch(
                        userDto -> Objects.equals(userDto.getId(), userToAdd1.getId().toString())
                )
        ).isTrue(); // check that user1 is in the circle
        Assertions.assertThat(
                circle.getUsers().stream().anyMatch(
                        userDto -> Objects.equals(userDto.getId(), userToAdd3.getId().toString())
                )
        ).isTrue(); // check that user3 is in the circle

        /* Delete circle */
        ResponseEntity<String> deleteCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(deleteCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Try to get deleted circle --> should return 404 */
        ResponseEntity<CircleDto> getCircleResponse5 = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleId),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        CircleDto.class
                );
        Assertions.assertThat(getCircleResponse5.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
    }

    @Test
    public void AddOrRemoveUserWhenUserIsNotTheCreator() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var wrongCircle = dummyDataCreator.generateCircle(true, null);
        var circleOfAuthenticatedUser = dummyDataCreator.generateCircle(true, authenticatedUser);
        var userToDoSomething = dummyDataCreator.generateUser(true);

        /* Try to add user in wrong circle --> should get 403 */
        ResponseEntity<String> addUserInWrongCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(wrongCircle.getId().toString())
                                .concat(HttpUtils.userUrl)
                                .concat(userToDoSomething.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addUserInWrongCircleResponse.getStatusCode()).isEqualTo(HttpStatus.FORBIDDEN);

        /* Try to add user in correct circle --> should get 200 */
        ResponseEntity<String> addUserInCorrectCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleOfAuthenticatedUser.getId().toString())
                                .concat(HttpUtils.userUrl)
                                .concat(userToDoSomething.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(addUserInCorrectCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Try to remove user from wrong circle --> should get 403 */
        ResponseEntity<String> removeUserFromWrongCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(wrongCircle.getId().toString())
                                .concat(HttpUtils.userUrl)
                                .concat(userToDoSomething.getId().toString()),
                        HttpMethod.DELETE,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeUserFromWrongCircleResponse.getStatusCode()).isEqualTo(HttpStatus.FORBIDDEN);

        /* Try to remove user from correct circle --> should get 200 */
        ResponseEntity<String> removeUserFromCorrectCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat(circleOfAuthenticatedUser.getId().toString())
                                .concat(HttpUtils.userUrl)
                                .concat(userToDoSomething.getId().toString()),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        String.class
                );
        Assertions.assertThat(removeUserFromCorrectCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
    }
}
