package com.teasy.CineCircleApi.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.UserFullInfoDto;
import com.teasy.CineCircleApi.models.dtos.requests.*;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.utils.SendEmailRequest;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.services.utils.EmailService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
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
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final EmailService emailService;

    private final static String resetPasswordUrlKey = "resetPasswordUrl";
    private final static String usernameKey = "username";
    private final static String resetPasswordMailSubject = "Réinitialisation de votre mot de passe";
    private final static String resetPasswordTemplateName = "reset-password.html";
    private final static String resetPasswordTUrlPrefix = "https://huco-reset-password.vercel.app/reset-pwd?token=";

    @Autowired
    public UserService(UserRepository userRepository,
                       PasswordEncoder passwordEncoder,
                       EmailService emailService) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
        this.emailService = emailService;
    }

    public UserFullInfoDto createUser(AuthSignUpRequest request) throws ExpectedException {
        // username should be only lowercase
        var finalUsername = request.username().toLowerCase();

        if (usernameAlreadyExists(finalUsername)) {
            throw new ExpectedException(ErrorDetails.ERR_USER_USERNAME_ALREADY_EXISTS.addingArgs(finalUsername));
        }

        if (emailAlreadyExists(request.email())) {
            throw new ExpectedException(ErrorDetails.ERR_USER_EMAIL_ALREADY_EXISTS.addingArgs(request.email()));
        }

        var user = new User(
                finalUsername,
                request.email(),
                passwordEncoder.encode(request.password()),
                request.displayName()
        );
        userRepository.save(user);

        return entityToFullInfoDto(user);
    }

    public UserFullInfoDto getUserFullInfo(String username) throws ExpectedException {
        return entityToFullInfoDto(findUserByUsernameOrElseThrow(username));
    }

    public UserDto resetPassword(String username, AuthResetPasswordRequest authResetPasswordRequest) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        // check if oldPassword i correct
        if (!passwordEncoder.matches(authResetPasswordRequest.oldPassword(), user.getHashPassword())) {
            throw new ExpectedException(ErrorDetails.ERR_USER_PASSWORD_INCORRECT.addingArgs(user.getId()));
        }

        // update password
        user.setHashPassword(passwordEncoder.encode(authResetPasswordRequest.newPassword()));
        userRepository.save(user);
        return entityToDto(user);
    }

    public UserDto resetPasswordWithToken(UserResetPasswordRequest userResetPasswordRequest) throws ExpectedException {
        var user = findUserByEmailOrElseThrow(userResetPasswordRequest.email());
        // check if token is correct
        if (!Objects.equals(user.getResetPasswordToken(), userResetPasswordRequest.token())) {
            throw new ExpectedException(ErrorDetails.ERR_USER_RESET_PASSWORD_TOKEN_INCORRECT.addingArgs(user.getId()));
        }

        // update password
        user.setHashPassword(passwordEncoder.encode(userResetPasswordRequest.newPassword()));
        user.setResetPasswordToken(null);
        userRepository.save(user);
        return entityToDto(user);
    }

    public UserFullInfoDto getUserFullInfoByUsernameOrEmail(String username, String email) throws ExpectedException {
        return entityToFullInfoDto(findUserByUsernameOrEmailOrElseThrow(username, email));
    }

    public UserDto getUser(UUID id) throws ExpectedException {
        return entityToDto(findUserByIdOrElseThrow(id));
    }

    public UserFullInfoDto updateUser(AuthMeUpdateRequest request, String username) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        user.setDisplayName(request.displayName());
        userRepository.save(user);
        return entityToFullInfoDto(user);
    }

    public Page<UserDto> searchUsers(String username, Pageable pageable, UserSearchRequest request) throws ExpectedException {
        // check query content
        if (request.query().isEmpty()) {
            throw new ExpectedException(ErrorDetails.ERR_USER_SEARCH_QUERY_EMPTY);
        }
        User user = findUserByUsernameOrElseThrow(username);
        // create example with query that can match username, email or displayName
        return userRepository.searchUsers(request.query(), user.getId(), pageable)
                .map(this::entityToDto);
    }

    public void sendResetPasswordEmail(String email) throws ExpectedException {
        var result = userRepository.findByEmail(email);
        // if user cannot be found, we should not let requester know it, it will avoid anyone knowing that an email exists in database
        if (result.isEmpty()) {
            log.warn("User cannot be found with email {}", email);
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
        templateValues.put(resetPasswordUrlKey, resetPasswordTUrlPrefix + token);
        SendEmailRequest sendEmailRequest = new SendEmailRequest(
                resetPasswordMailSubject,
                email,
                resetPasswordTemplateName,
                templateValues);
        emailService.sendEmail(sendEmailRequest);
    }

    public UserFullInfoDto addUserToRelatedUsers(String username, UUID relatedUserId) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        if (user.getId().equals(relatedUserId)) {
            return entityToFullInfoDto(user);
        }
        var relatedUser = findUserByIdOrElseThrow(relatedUserId);
        user.addRelatedUser(relatedUser);
        userRepository.save(user);
        return entityToFullInfoDto(user);
    }

    public UserFullInfoDto removeUserFromRelatedUsers(String username, UUID relatedUserId) throws ExpectedException {
        var user = findUserByUsernameOrElseThrow(username);
        var relatedUser = findUserByIdOrElseThrow(relatedUserId);
        user.removeRelatedUser(relatedUser);
        userRepository.save(user);
        return entityToFullInfoDto(user);
    }

    public User findUserByUsernameOrElseThrow(String username) throws ExpectedException {
        return userRepository
                .findByUsername(username)
                .orElseThrow(() -> new ExpectedException(
                        ErrorDetails.ERR_USER_NOT_FOUND.addingArgs(username))
                );
    }

    public User findUserByEmailOrElseThrow(String email) throws ExpectedException {
        return userRepository
                .findByEmail(email)
                .orElseThrow(() -> new ExpectedException(
                        ErrorDetails.ERR_USER_NOT_FOUND.addingArgs(email))
                );
    }

    public User findUserByUsernameOrEmailOrElseThrow(String username, String email) throws ExpectedException {
        return userRepository
                .findByUsernameOrEmail(username, email)
                .orElseThrow(() -> new ExpectedException(
                        ErrorDetails.ERR_USER_NOT_FOUND.addingArgs(String.format("%s / %s", username, email)))
                );
    }

    public User findUserByIdOrElseThrow(UUID id) throws ExpectedException {
        return userRepository
                .findById(id)
                .orElseThrow(() -> new ExpectedException(
                        ErrorDetails.ERR_USER_NOT_FOUND.addingArgs(id))
                );
    }

    public void save(User user) {
        userRepository.save(user);
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
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(user, UserDto.class);
    }

    private UserFullInfoDto entityToFullInfoDto(User user) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(user, UserFullInfoDto.class);
    }
}
