package com.teasy.CineCircleApi.services;

//import com.teasy.CineCircleApi.repositories.UserRepository;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.requests.*;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.entities.User;
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
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

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

    public UserDto createUser(AuthSignUpRequest request) throws ResponseStatusException {
        if (usernameAlreadyExists(request.username())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST,
                    String.format("username %s already exists", request.username()));
        }

        if (emailAlreadyExists(request.email())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST,
                    String.format("email %s already exists", request.email()));
        }

        var user = new User(request.username(),
                request.email(),
                passwordEncoder.encode(request.password()),
                request.displayName());
        userRepository.save(user);

        return entityToDto(user);
    }

    public UserDto getUserByUsername(String username) throws ResponseStatusException {
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with username %s", username)));
        return entityToDto(user);
    }

    public UserFullInfoDto getUserFullInfo(String username) throws ResponseStatusException {
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with username %s", username)));
        return entityToFullInfoDto(user);
    }

    public UserDto resetPassword(String username, AuthResetPasswordRequest authResetPasswordRequest) throws ResponseStatusException {
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with username %s", username)));
        // check if oldPassword i correct
        if (!passwordEncoder.matches(authResetPasswordRequest.oldPassword(), user.getHashPassword())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN,
                    String.format("oldPassword for user with username %s is incorrect", username));
        }

        // update password
        user.setHashPassword(passwordEncoder.encode(authResetPasswordRequest.newPassword()));
        userRepository.save(user);
        return entityToDto(user);
    }

    public UserDto resetPasswordWithToken(UserResetPasswordRequest userResetPasswordRequest) throws ResponseStatusException {
        var user = userRepository
                .findByEmail(userResetPasswordRequest.email())
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with email %s", userResetPasswordRequest.email())));
        // check if token ius correct
        if (!Objects.equals(user.getResetPasswordToken(), userResetPasswordRequest.token())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN,
                    String.format("token for user with email %s is incorrect", userResetPasswordRequest.email()));
        }

        // update password
        user.setHashPassword(passwordEncoder.encode(userResetPasswordRequest.newPassword()));
        user.setResetPasswordToken(null);
        userRepository.save(user);
        return entityToDto(user);
    }

    public UserFullInfoDto getUserByUsernameOrEmail(String username, String email) throws ResponseStatusException {
        var user = userRepository
                .findByUsernameOrEmail(username, email)
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with username %s or email %s", username, email)));
        return entityToFullInfoDto(user);
    }

    public UserDto getUser(Long id) throws ResponseStatusException {
        var user = userRepository
                .findById(id)
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with id %d", id)));
        return entityToDto(user);
    }

    public UserDto updateUser(AuthMeUpdateRequest request, String username) throws ResponseStatusException {
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with username %s", username)));
        user.setDisplayName(request.getDisplayName());
        user = userRepository.save(user);
        return entityToDto(user);
    }

    public Page<UserDto> searchUsers(Pageable pageable, UserSearchRequest request) throws ResponseStatusException {
        // check query content
        if (request.query().isEmpty()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "query cannot be empty");
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
