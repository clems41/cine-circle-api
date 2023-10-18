package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.responses.CustomErrorResponse;
import org.springframework.http.ResponseEntity;
import org.springframework.web.server.ResponseStatusException;

public class HttpErrorService {
    public static ResponseEntity<CustomErrorResponse> getEntityResponseFromException(ResponseStatusException e) {
        return new ResponseEntity<CustomErrorResponse>(new CustomErrorResponse(e), e.getStatusCode());
    }
}
