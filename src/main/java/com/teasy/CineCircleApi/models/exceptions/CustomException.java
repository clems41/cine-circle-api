package com.teasy.CineCircleApi.models.exceptions;

import com.teasy.CineCircleApi.models.dtos.responses.CustomErrorResponse;
import com.teasy.CineCircleApi.models.enums.ErrorCode;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import org.springframework.web.server.ResponseStatusException;

public class CustomException extends ResponseStatusException {

    public CustomException(HttpStatusCode status, ErrorCode errorCode, String message) {
        super(status, errorCode.name(), new Throwable(message));
    }

    public ResponseEntity<CustomErrorResponse> getEntityResponse() {
        return new ResponseEntity<>(new CustomErrorResponse(this), this.getStatusCode());
    }
}
