package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.dtos.CircleDto;
import com.teasy.CineCircleApi.models.dtos.CirclePublicDto;
import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.entities.Circle;
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

import java.util.*;

public class CircleTest extends IntegrationTestAbstract {
    @Test
    public void CRUDCircleAddAndRemoveUsers() {
        /* Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();
        var userToAdd1 = dummyDataCreator.generateUser(true);
        var userToAdd2 = dummyDataCreator.generateUser(true);
        var userToAdd3 = dummyDataCreator.generateUser(true);
        var name1 = RandomUtils.randomString(15);
        var name2 = RandomUtils.randomString(15);
        var description1 = RandomUtils.randomString(30);
        var description2 = RandomUtils.randomString(30);

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
        var wrongCircle = dummyDataCreator.generateCircle(true, null, null, null);
        var circleOfAuthenticatedUser = dummyDataCreator.generateCircle(true, authenticatedUser, null, null);
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

    @Test
    public void searchPublicCircles() {
        /*  Data */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        List<Circle> matchingCircles = new ArrayList<>();
        var keyword = RandomUtils.randomString(5);

        /* Create some non-matching keyword private circles */
        for (int i = 0; i < RandomUtils.randomInt(5, 10); i++) {
            dummyDataCreator.generateCircle(true, null, false, null);
        }

        /* Create some non-matching keyword public circles */
        for (int i = 0; i < RandomUtils.randomInt(5, 10); i++) {
            dummyDataCreator.generateCircle(true, null, true, null);
        }

        /* Create some matching keyword private circles */
        for (int i = 0; i < RandomUtils.randomInt(2, 5); i++) { // with keyword at the beginning
            var matchingName = keyword.concat(RandomUtils.randomString(15));
            dummyDataCreator.generateCircle(true, null, false, matchingName);
        }
        for (int i = 0; i < RandomUtils.randomInt(2, 5); i++) { // with keyword at the end
            var matchingName = RandomUtils.randomString(15).concat(keyword);
            dummyDataCreator.generateCircle(true, null, false, matchingName);
        }
        for (int i = 0; i < RandomUtils.randomInt(2, 5); i++) { // with keyword in the middle
            var matchingName = RandomUtils.randomString(7)
                    .concat(keyword)
                    .concat(RandomUtils.randomString(7));
            dummyDataCreator.generateCircle(true, null, false, matchingName);
        }

        /* Create some matching keyword public circles */
        for (int i = 0; i < RandomUtils.randomInt(2, 5); i++) { // with keyword at the beginning
            var matchingName = keyword.concat(RandomUtils.randomString(15));
            matchingCircles.add(dummyDataCreator.generateCircle(true, null, true, matchingName));
        }
        for (int i = 0; i < RandomUtils.randomInt(2, 5); i++) { // with keyword at the end
            var matchingName = RandomUtils.randomString(15).concat(keyword);
            matchingCircles.add(dummyDataCreator.generateCircle(true, null, true, matchingName));
        }
        for (int i = 0; i < RandomUtils.randomInt(2, 5); i++) { // with keyword in the middle
            var matchingName = RandomUtils.randomString(7)
                    .concat(keyword)
                    .concat(RandomUtils.randomString(7));
            matchingCircles.add(dummyDataCreator.generateCircle(true, null, true, matchingName));
        }

        /* Search public with keyword */
        Map<String, Object> queryParams = new HashMap<>();
        queryParams.put("page", 0);
        queryParams.put("size", 15);
        queryParams.put("query", keyword);
        ResponseEntity<CustomPageImpl<CirclePublicDto>> searchPublicCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getUriWithQueryParameter(port, HttpUtils.circleUrl.concat("public"), queryParams),
                        HttpMethod.GET,
                        new HttpEntity<>(null, headers),
                        new ParameterizedTypeReference<>(){}
                );
        Assertions.assertThat(searchPublicCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);

        /* Check that response contains only public matching keyword circles */
        Assertions.assertThat(searchPublicCircleResponse.getBody()).isNotNull();
        var actualCircles = searchPublicCircleResponse.getBody().stream().toList();
        Assertions.assertThat(matchingCircles.size()).isEqualTo(actualCircles.size());
        actualCircles.forEach(actualCircle -> {
            // find if circle is one of the expected
            var matchingCircle = matchingCircles
                    .stream()
                    .filter(circle -> Objects.equals(circle.getId().toString(), actualCircle.getId()))
                    .findAny();
            Assertions.assertThat(matchingCircle.isPresent()).isTrue();
            var expectedCircle = matchingCircle.get();
            // compare fields value
            Assertions.assertThat(actualCircle.getIsPublic()).isEqualTo(expectedCircle.getIsPublic());
            Assertions.assertThat(actualCircle.getName()).isEqualTo(expectedCircle.getName());
            Assertions.assertThat(actualCircle.getDescription()).isEqualTo(expectedCircle.getDescription());
            Assertions.assertThat(actualCircle.getCreatedBy().getId()).isEqualTo(expectedCircle.getCreatedBy().getId().toString());
        });
    }

    @Test
    public void JoinPublicCircle() {
        /* Create circle */
        var circle = dummyDataCreator.generateCircle(true, null, true, null);

        /* Create user */
        var signUpRequest = authenticator.authenticateNewUser();
        var headers = authenticator.authenticateUserAndGetHeadersWithJwtToken(signUpRequest.username(), signUpRequest.password());
        var authenticatedUser = userRepository.findByUsername(signUpRequest.username()).orElseThrow();

        /* User joins circle */
        ResponseEntity<CirclePublicDto> joinPublicCircleResponse = this.restTemplate
                .exchange(
                        HttpUtils.getTestingUrl(port)
                                .concat(HttpUtils.circleUrl)
                                .concat("public/")
                                .concat(circle.getId().toString())
                                .concat("/join"),
                        HttpMethod.PUT,
                        new HttpEntity<>(null, headers),
                        CirclePublicDto.class
                );
        Assertions.assertThat(joinPublicCircleResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        Assertions.assertThat(joinPublicCircleResponse.getBody()).isNotNull();
        Assertions.assertThat(joinPublicCircleResponse.getBody().getId()).isEqualTo(circle.getId().toString());

        /* Check that circle contains user in database */
        var updatedCircle = circleRepository.findById(circle.getId()).orElseThrow();
        Assertions.assertThat(updatedCircle.getUsers().size()).isEqualTo(circle.getUsers().size() + 1);
        Assertions.assertThat(updatedCircle.getUsers().stream().anyMatch(user -> user.getId().equals(authenticatedUser.getId()))).isTrue();
    }
}
