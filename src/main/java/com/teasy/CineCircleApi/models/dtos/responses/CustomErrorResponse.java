package com.teasy.CineCircleApi.models.dtos.responses;

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
    private String errorStack;

    public CustomErrorResponse(ResponseStatusException e) {
        if (e != null) {
            this.errorMessage = e.getMessage();
            this.errorStack = Arrays.toString(e.getStackTrace());
        }
    }
}
