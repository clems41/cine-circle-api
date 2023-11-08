package com.teasy.CineCircleApi.models.utils;

import jakarta.validation.constraints.Email;

import java.util.Map;

public record SendEmailRequest(
        String subject,
        @Email(regexp = ".+@.+\\..+") String receiver,
        String templateName,
        Map<String, String> templateValues
) {
}
