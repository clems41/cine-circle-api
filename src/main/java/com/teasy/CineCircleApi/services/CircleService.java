package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.CircleDto;
import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.repositories.CircleRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

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
        var user = findUserByUsernameOrElseThrow(authenticatedUsername);
        var circles = circleRepository.findAllByUsers_Id(user.getId());
        return circles
                .stream()
                .map(this::fromCircleEntityToCircleDto)
                .toList();
    }

    public CircleDto createCircle(CircleCreateUpdateRequest circleCreateUpdateRequest, String authenticatedUsername) {
        var user = findUserByUsernameOrElseThrow(authenticatedUsername);
        var newCircle = new Circle(
                circleCreateUpdateRequest.isPublic(),
                circleCreateUpdateRequest.name(),
                circleCreateUpdateRequest.description(),
                user
        );
        circleRepository.save(newCircle);
        return fromCircleEntityToCircleDto(newCircle);
    }

    public CircleDto updateCircle(CircleCreateUpdateRequest circleCreateUpdateRequest, Long circleId, String authenticatedUsername) throws CustomException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // update circle
        circle.setDescription(circleCreateUpdateRequest.description());
        circle.setName(circleCreateUpdateRequest.name());
        circle.setIsPublic(circleCreateUpdateRequest.isPublic());
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle);
    }

    public void deleteCircle(Long circleId, String authenticatedUsername) throws CustomException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // delete circle
        circleRepository.delete(circle);
    }

    public CircleDto addUserToCircle(Long userIdToAdd, Long circleId, String authenticatedUsername) throws CustomException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // find user to add
        var userToAdd = findUserByIdOrElseThrow(userIdToAdd);

        // add user
        circle.addUser(userToAdd);
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle);
    }

    public CircleDto removeUserFromCircle(Long userIdToRemove, Long circleId, String authenticatedUsername) throws CustomException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // find user to remove
        var userToRemove = findUserByIdOrElseThrow(userIdToRemove);

        // add user
        circle.removeUser(userToRemove);
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle);
    }

    private Circle getCircleAndCheckPermissions(Long circleId, String authenticatedUsername) {
        var user = findUserByUsernameOrElseThrow(authenticatedUsername);
        var circle = findCircleByIdOrElseThrow(circleId);

        // user should the one created the circle to update or delete it
        if (!Objects.equals(circle.getCreatedBy().getId(), user.getId())) {
            throw CustomExceptionHandler.circleWithIdUserWithUsernameBadPermissions(circleId, authenticatedUsername);
        }

        return circle;
    }

    private User findUserByIdOrElseThrow(Long id) throws CustomException {
        return userRepository
                .findById(id)
                .orElseThrow(() -> CustomExceptionHandler.userWithIdNotFound(id));
    }

    private User findUserByUsernameOrElseThrow(String username) throws CustomException {
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
    }

    private Circle findCircleByIdOrElseThrow(Long circleId) {
        // check if circle exist
        return circleRepository
                .findById(circleId)
                .orElseThrow(() -> CustomExceptionHandler.circleWithIdNotFound(circleId));
    }

    private CircleDto fromCircleEntityToCircleDto(Circle circle) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(circle, CircleDto.class);
    }
}
