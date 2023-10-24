package com.teasy.CineCircleApi.models.dtos.responses;

import com.teasy.CineCircleApi.models.exceptions.CustomException;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import org.springframework.web.server.ResponseStatusException;

import java.util.Arrays;

@Getter
public class CustomErrorResponse {
    private String errorMessage;
    private String errorCode;
    private StackTraceElement[] errorStack;

    public CustomErrorResponse(CustomException e) {
        if (e != null) {
            this.errorCode = e.getReason();
            this.errorMessage = e.getCause().getMessage();
            this.errorStack = e.getStackTrace();
        }
    }
}
