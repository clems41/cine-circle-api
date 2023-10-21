package com.teasy.CineCircleApi.services;

//import com.teasy.CineCircleApi.repositories.UserRepository;

import com.teasy.CineCircleApi.models.dtos.UserDto;
import com.teasy.CineCircleApi.models.dtos.requests.SignUpRequest;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

@Service
public class UserService {
    UserRepository userRepository;
    PasswordEncoder passwordEncoder;

    @Autowired
    public UserService(UserRepository userRepository,
                       PasswordEncoder passwordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    public UserDto createUser(SignUpRequest request) throws ResponseStatusException {
        if (usernameAlreadyExists(request.username())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST,
                    String.format("username %s already exists", request.username()));
        }

        if (emailAlreadyExists(request.email())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST,
                    String.format("email %s already exists", request.email()));
        }

        var user = new User(request.username(), request.email(), passwordEncoder.encode(request.password()));
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

    public UserDto getUserByUsernameOrEmail(String username, String email) throws ResponseStatusException {
        var user = userRepository
                .findByUsernameOrEmail(username, email)
                .orElseThrow(() ->
                        new ResponseStatusException(HttpStatus.NOT_FOUND,
                                String.format("user cannot be found with username %s or email %s", username, email)));
        return entityToDto(user);
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
        var dto = new UserDto();
        if (user != null) {
            dto.setUsername(user.getUsername());
            dto.setEmail(user.getEmail());
            dto.setId(user.getId());
        }
        return dto;
    }
}
