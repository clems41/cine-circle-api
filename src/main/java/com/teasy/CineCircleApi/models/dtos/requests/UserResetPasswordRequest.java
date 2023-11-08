package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotEmpty;
import lombok.NonNull;
import org.hibernate.validator.constraints.Length;

public record UserResetPasswordRequest(
        @NotEmpty @Email(regexp = ".+@.+\\..+") String email,
        @NotEmpty String token,
        @NotEmpty @Length(min = 6) String newPassword
) {
}
