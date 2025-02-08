package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;

public record ContactSendFeedbackRequest(
    @NotEmpty String feedback
){
}
