package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.CircleDto;
import com.teasy.CineCircleApi.models.dtos.CirclePublicDto;
import com.teasy.CineCircleApi.models.dtos.requests.CircleCreateUpdateRequest;
import com.teasy.CineCircleApi.models.dtos.requests.CircleSearchPublicRequest;
import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.CircleRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;
import java.util.UUID;

@Service
public class CircleService {
    private final CircleRepository circleRepository;
    private final UserService userService;

    @Autowired
    public CircleService(CircleRepository circleRepository, UserService userService) {
        this.circleRepository = circleRepository;
        this.userService = userService;
    }

    public List<CircleDto> listCircles(String authenticatedUsername) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
        var circles = circleRepository.findAllByUsers_Id(user.getId());
        return circles
                .stream()
                .map(circle -> fromCircleEntityToCircleDto(circle, CircleDto.class))
                .toList();
    }

    public CirclePublicDto joinPublicCircle(UUID circleId, String authenticatedUsername) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
        var circle = findCircleByIdOrElseThrow(circleId);
        if (!circle.getIsPublic()) {
            throw new ExpectedException(ErrorDetails.ERR_CIRCLE_USER_BAD_PERMISSIONS.addingArgs(authenticatedUsername));
        }
        circle.addUser(user);
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle, CirclePublicDto.class);
    }

    public Page<CirclePublicDto> searchPublicCircles(CircleSearchPublicRequest circleSearchPublicRequest, Pageable pageable) {
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues()
                .withStringMatcher(ExampleMatcher.StringMatcher.CONTAINING);
        var matchingCircle = new Circle();
        matchingCircle.setName(circleSearchPublicRequest.query());
        matchingCircle.setIsPublic(true);

        var circles = circleRepository.findAll(Example.of(matchingCircle, matcher), pageable);
        return circles.map(circle -> fromCircleEntityToCircleDto(circle, CirclePublicDto.class));
    }

    public CircleDto createCircle(CircleCreateUpdateRequest circleCreateUpdateRequest, String authenticatedUsername) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
        var newCircle = new Circle(
                circleCreateUpdateRequest.isPublic(),
                circleCreateUpdateRequest.name(),
                circleCreateUpdateRequest.description(),
                user
        );
        circleRepository.save(newCircle);
        return fromCircleEntityToCircleDto(newCircle, CircleDto.class);
    }

    public CircleDto updateCircle(CircleCreateUpdateRequest circleCreateUpdateRequest, UUID circleId, String authenticatedUsername) throws ExpectedException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // update circle
        circle.setDescription(circleCreateUpdateRequest.description());
        circle.setName(circleCreateUpdateRequest.name());
        circle.setIsPublic(circleCreateUpdateRequest.isPublic());
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle, CircleDto.class);
    }

    public void deleteCircle(UUID circleId, String authenticatedUsername) throws ExpectedException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // delete circle
        circleRepository.delete(circle);
    }

    public CircleDto addUserToCircle(UUID userIdToAdd, UUID circleId, String authenticatedUsername) throws ExpectedException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // find user to add
        var userToAdd = userService.findUserByIdOrElseThrow(userIdToAdd);

        // add user
        circle.addUser(userToAdd);
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle, CircleDto.class);
    }

    public CircleDto removeUserFromCircle(UUID userIdToRemove, UUID circleId, String authenticatedUsername) throws ExpectedException {
        var circle = getCircleAndCheckPermissions(circleId, authenticatedUsername);

        // find user to remove
        var userToRemove = userService.findUserByIdOrElseThrow(userIdToRemove);

        // add user
        circle.removeUser(userToRemove);
        circleRepository.save(circle);
        return fromCircleEntityToCircleDto(circle, CircleDto.class);
    }

    public CircleDto getCircle(UUID circleId, String authenticatedUsername) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
        var circle = findCircleByIdOrElseThrow(circleId);

        // if user is not in circle he cannot get circle info
        if (circle.getUsers().stream().noneMatch(circleUser -> circleUser.getId() == user.getId())) {
            throw new ExpectedException(ErrorDetails.ERR_CIRCLE_USER_BAD_PERMISSIONS.addingArgs(authenticatedUsername));
        }
        return fromCircleEntityToCircleDto(circle, CircleDto.class);
    }

    public Circle findCircleByIdOrElseThrow(UUID circleId) throws ExpectedException {
        // check if circle exist
        return circleRepository
                .findById(circleId)
                .orElseThrow(() -> new ExpectedException(ErrorDetails.ERR_CIRCLE_NOT_FOUND.addingArgs(circleId)));
    }

    private Circle getCircleAndCheckPermissions(UUID circleId, String authenticatedUsername) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(authenticatedUsername);
        var circle = findCircleByIdOrElseThrow(circleId);

        // user should the one created the circle to update or delete it
        if (!Objects.equals(circle.getCreatedBy().getId(), user.getId())) {
            throw new ExpectedException(ErrorDetails.ERR_CIRCLE_USER_BAD_PERMISSIONS.addingArgs(authenticatedUsername));
        }

        return circle;
    }

    private <T> T fromCircleEntityToCircleDto(Circle circle, Class<T> circleDtoType) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(circle, circleDtoType);
    }
}
