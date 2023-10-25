package com.teasy.CineCircleApi.services;

//import com.teasy.CineCircleApi.repositories.UserRepository;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.*;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.models.utils.SendEmailRequest;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.services.utils.EmailService;
import jakarta.mail.MessagingException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.UUID;

@Service
@Slf4j
public class UserService {
    UserRepository userRepository;
    PasswordEncoder passwordEncoder;
    EmailService emailService;

    private final static String resetPasswordUrlKey = "resetPasswordUrl";
    private final static String usernameKey = "username";
    private final static String tokenKey = "token";
    private final static String resetPasswordMailSubject = "Réinitialisation de votre mot de passe";
    private final static String resetPasswordTemplateName = "reset-password.html";

    @Autowired
    public UserService(UserRepository userRepository,
                       PasswordEncoder passwordEncoder,
                       EmailService emailService) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
        this.emailService = emailService;
    }

    public UserDto createUser(AuthSignUpRequest request) throws CustomException {
        if (usernameAlreadyExists(request.username())) {
            throw CustomExceptionHandler.userWithUsernameAlreadyExists(request.username());
        }

        if (emailAlreadyExists(request.email())) {
            throw CustomExceptionHandler.userWithEmailAlreadyExists(request.email());
        }

        var user = new User(
                request.username(),
                request.email(),
                passwordEncoder.encode(request.password()),
                request.displayName()
        );
        userRepository.save(user);

        return entityToDto(user);
    }

    public UserFullInfoDto getUserFullInfo(String username) throws CustomException {
        return entityToFullInfoDto(getUserWithUsernameOrElseThrow(username));
    }

    public UserDto resetPassword(String username, AuthResetPasswordRequest authResetPasswordRequest) throws CustomException {
        var user = getUserWithUsernameOrElseThrow(username);
        // check if oldPassword i correct
        if (!passwordEncoder.matches(authResetPasswordRequest.oldPassword(), user.getHashPassword())) {
            throw CustomExceptionHandler.userWithUsernameBadCredentials(username);
        }

        // update password
        user.setHashPassword(passwordEncoder.encode(authResetPasswordRequest.newPassword()));
        userRepository.save(user);
        return entityToDto(user);
    }

    public UserDto resetPasswordWithToken(UserResetPasswordRequest userResetPasswordRequest) throws CustomException {
        var user = getUserWithEmailOrElseThrow(userResetPasswordRequest.email());
        // check if token ius correct
        if (!Objects.equals(user.getResetPasswordToken(), userResetPasswordRequest.token())) {
            throw CustomExceptionHandler.userWithEmailBadCredentials(userResetPasswordRequest.email());
        }

        // update password
        user.setHashPassword(passwordEncoder.encode(userResetPasswordRequest.newPassword()));
        user.setResetPasswordToken(null);
        userRepository.save(user);
        return entityToDto(user);
    }

    public UserFullInfoDto getUserFullInfoByUsernameOrEmail(String username, String email) throws CustomException {
        return entityToFullInfoDto(getUserWithUsernameOrEmailOrElseThrow(username, email));
    }

    public UserDto getUser(Long id) throws CustomException {
        return entityToDto(getUserWithIdOrElseThrow(id));
    }

    public UserDto updateUser(AuthMeUpdateRequest request, String username) throws CustomException {
        var user = getUserWithUsernameOrElseThrow(username);
        user.setDisplayName(request.getDisplayName());
        userRepository.save(user);
        return entityToDto(user);
    }

    public Page<UserDto> searchUsers(Pageable pageable, UserSearchRequest request) throws CustomException {
        // check query content
        if (request.query().isEmpty()) {
            throw CustomExceptionHandler.userSearchQueryEmpty();
        }
        // create example with query that can match username, email or displayName
        ExampleMatcher matcher = ExampleMatcher
                .matchingAny()
                .withStringMatcher(ExampleMatcher.StringMatcher.CONTAINING)
                .withIgnoreCase()
                .withIgnoreNullValues();
        var exampleUser = new User();
        exampleUser.setUsername(request.query());
        exampleUser.setDisplayName(request.query());

        // find users
        var users = userRepository
                .findAll(Example.of(exampleUser, matcher), pageable);
        return users.map(this::entityToDto);
    }

    public void sendResetPasswordEmail(String email) {
        var result = userRepository.findByEmail(email);
        // if user cannot be found, we should not let requester know it, it will avoid anyone knowing that an email exists in database
        if (result.isEmpty()) {
            return;
        }
        var user = result.get();

        // generate token that will be used to reset password
        var token = UUID.randomUUID().toString();
        user.setResetPasswordToken(token);
        userRepository.save(user);

        // send email if user exists
        Map<String, String> templateValues = new HashMap<>();
        templateValues.put(usernameKey, user.getUsername());
        templateValues.put(resetPasswordUrlKey, "TODO");
        templateValues.put(tokenKey, token);
        SendEmailRequest sendEmailRequest = new SendEmailRequest(
                resetPasswordMailSubject,
                email,
                resetPasswordTemplateName,
                templateValues);
        try {
            emailService.sendEmail(sendEmailRequest);
        } catch (MessagingException e) {
            log.error("cannot send reset password email : {}", e.getMessage());
        }
    }

    private User getUserWithUsernameOrElseThrow(String username) throws CustomException {
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
    }

    private User getUserWithEmailOrElseThrow(String email) throws CustomException {
        return userRepository
                .findByEmail(email)
                .orElseThrow(() -> CustomExceptionHandler.userWithEmailNotFound(email));
    }

    private User getUserWithUsernameOrEmailOrElseThrow(String username, String email) throws CustomException {
        return userRepository
                .findByUsernameOrEmail(username, email)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameOrEmailNotFound(username, email));
    }

    private User getUserWithIdOrElseThrow(Long id) throws CustomException {
        return userRepository
                .findById(id)
                .orElseThrow(() -> CustomExceptionHandler.userWithIdNotFound(id));
    }

    private Boolean usernameAlreadyExists(String username) {
        var user = userRepository.findByUsername(username);
        return user.isPresent();
    }

    private Boolean emailAlreadyExists(String email) {
        var user = userRepository.findByEmail(email);
        return user.isPresent();
    }

    private UserDto entityToDto(User user) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        return mapper.convertValue(user, UserDto.class);
    }

    private UserFullInfoDto entityToFullInfoDto(User user) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        return mapper.convertValue(user, UserFullInfoDto.class);
    }
}
