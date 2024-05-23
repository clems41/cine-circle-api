package com.teasy.CineCircleApi.config;

import com.teasy.CineCircleApi.models.exceptions.ErrorResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

@Order(Ordered.HIGHEST_PRECEDENCE)
@ControllerAdvice
public class RestExceptionHandler extends ResponseEntityExceptionHandler {


    @ExceptionHandler(ExpectedException.class)
    // when exception is expected, we should use its HttpStatus when returning ErrorResponse
    private ResponseEntity<ErrorResponse> buildResponseEntityForExpectedException(ExpectedException exception) {
        return new ResponseEntity<>(buildErrorResponse(exception, exception.getError().getCode()), exception.getHttpStatus());
    }

    @ExceptionHandler(Exception.class)
    // when exception is not expected, we should use HttpStatus.INTERNAL_SERVER_ERROR when returning ErrorResponse
    private ResponseEntity<ErrorResponse> buildResponseEntityForException(Exception exception) {
        return new ResponseEntity<>(buildErrorResponse(exception, null), HttpStatus.INTERNAL_SERVER_ERROR);
    }

    private ErrorResponse buildErrorResponse(Exception exception, String errorCode) {
        return new ErrorResponse(
                exception.getMessage(),
                errorCode,
                exception.getCause() != null ? exception.getCause().getMessage() : "",
                exception.getCause() != null ? exception.getCause().getStackTrace() : null
        );
    }
}
