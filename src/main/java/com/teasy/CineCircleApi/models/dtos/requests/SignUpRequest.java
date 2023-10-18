package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.Email;
import org.hibernate.validator.constraints.Length;
import lombok.NonNull;

public record SignUpRequest(
        @NonNull String username,
        @NonNull @Email(regexp = ".+[@].+[\\.].+") String email,
        @NonNull @Length(min = 6) String password
) {}
