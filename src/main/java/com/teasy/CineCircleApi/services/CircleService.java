package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.CircleDto;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.repositories.CircleRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.util.List;
import java.util.Objects;

@Service
public class CircleService {
    CircleRepository circleRepository;
    UserRepository userRepository;

    @Autowired
    public CircleService(CircleRepository circleRepository, UserRepository userRepository) {
        this.circleRepository = circleRepository;
        this.userRepository = userRepository;
    }

    public List<CircleDto> listCircles(String authenticatedUsername) {
        var user = getUserAndThrownIfNotExist(authenticatedUsername);
        var circles = circleRepository.findAllByUsers_Id(user.getId());
        return circles
                .stream()
                .map(this::fromCircleEntityToCircleDto)
                .toList();
    }

    public CircleDto createCircle(CircleCreateUpdateRequest circleCreateUpdateRequest, String authenticatedUsername) {
        var user = getUserAndThrownIfNotExist(authenticatedUsername);
        var newCircle = new Circle(circleCreateUpdateRequest.isPublic(),
                circleCreateUpdateRequest.name(),
                circleCreateUpdateRequest.description(),
                user);
        circleRepository.save(newCircle);
        return fromCircleEntityToCircleDto(newCircle);
    }

    public CircleDto updateCircle(CircleCreateUpdateRequest circleCreateUpdateRequest, Long circleId, String authenticatedUsername) throws ResponseStatusException {
        var user = getUserAndThrownIfNotExist(authenticatedUsername);
        var circle = getCircleAndThrownIfNotExist(circleId);

        // user should the one created the circle to update it
        if (!Objects.equals(circle.getCreatedBy().getId(), user.getId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN,
                    String.format("user with id %d cannot update the circle with id %d because he did not created it",
                            user.getId(), circleId));
        }

        // update circle
        circle.setDescription(circleCreateUpdateRequest.description());
        circle.setName(circleCreateUpdateRequest.name());
        circle.setIsPublic(circleCreateUpdateRequest.isPublic());
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle);
    }

    public void deleteCircle(Long circleId, String authenticatedUsername) throws ResponseStatusException {
        var user = getUserAndThrownIfNotExist(authenticatedUsername);
        var circle = getCircleAndThrownIfNotExist(circleId);

        // user should the one created the circle to delete it
        if (!Objects.equals(circle.getCreatedBy().getId(), user.getId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN,
                    String.format("user with id %d cannot delete the circle with id %d because he did not created it",
                            user.getId(), circleId));
        }

        // delete circle
        circleRepository.delete(circle);
    }

    public CircleDto addUserToCircle(Long userIdToAdd, Long circleId, String authenticatedUsername) throws ResponseStatusException {
        var user = getUserAndThrownIfNotExist(authenticatedUsername);
        var circle = getCircleAndThrownIfNotExist(circleId);

        // user should the one created the circle to update it
        if (!Objects.equals(circle.getCreatedBy().getId(), user.getId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN,
                    String.format("user with id %d cannot update the circle with id %d because he did not created it",
                            user.getId(), circleId));
        }

        // find user to add
        var userToAdd = userRepository.findById(userIdToAdd).orElseThrow(() -> new ResponseStatusException(
                HttpStatus.NOT_FOUND,
                String.format("user with id %d cannot be found", userIdToAdd)));

        // add user
        circle.addUser(userToAdd);
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle);
    }

    public CircleDto removeUserFromCircle(Long userIdToRemove, Long circleId, String authenticatedUsername) throws ResponseStatusException {
        var user = getUserAndThrownIfNotExist(authenticatedUsername);
        var circle = getCircleAndThrownIfNotExist(circleId);

        // user should the one created the circle to update it
        if (!Objects.equals(circle.getCreatedBy().getId(), user.getId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN,
                    String.format("user with id %d cannot update the circle with id %d because he did not created it",
                            user.getId(), circleId));
        }

        // find user to remove
        var userToRemove = userRepository.findById(userIdToRemove).orElseThrow(() -> new ResponseStatusException(
                HttpStatus.NOT_FOUND,
                String.format("user with id %d cannot be found", userIdToRemove)));

        // add user
        circle.removeUser(userToRemove);
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle);
    }

    private User getUserAndThrownIfNotExist(String username) {
        // check if user exist
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("user with username %s cannot be found", username)));
    }

    private Circle getCircleAndThrownIfNotExist(Long circleId) {
        // check if circle exist
        return circleRepository
                .findById(circleId)
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("circle with id %d cannot be found", circleId)));
    }

    private CircleDto fromCircleEntityToCircleDto(Circle circle) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(circle, CircleDto.class);
    }
}
