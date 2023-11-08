package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotEmpty;
import org.hibernate.validator.constraints.Length;
import lombok.NonNull;

public record AuthSignUpRequest(
        @NotEmpty @Length(min = 6) String username,
        @NotEmpty @Email(regexp = ".+@.+\\..+") String email,
        @NotEmpty @Length(min = 6) String password,
        @NotEmpty @Length(min = 4, max = 20) String displayName
) {}
