package com.teasy.CineCircleApi.models.dtos.requests;

import lombok.NonNull;
import org.springframework.web.bind.annotation.RequestParam;


public record UserSendResetPasswordEmailRequest(
    @RequestParam(defaultValue = "") @NonNull String email
){
}
