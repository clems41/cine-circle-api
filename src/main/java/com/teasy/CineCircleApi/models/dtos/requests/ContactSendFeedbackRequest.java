package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;

public record ContactSendFeedbackRequest(
    @NotEmpty(message = "ERR_CONTACT_FEEDBACK_EMPTY") String feedback
){
}
