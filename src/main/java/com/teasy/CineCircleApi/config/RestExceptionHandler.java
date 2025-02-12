package com.teasy.CineCircleApi.config;

import com.teasy.CineCircleApi.models.entities.Error;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ErrorResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.ErrorRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

@Order(Ordered.HIGHEST_PRECEDENCE)
@ControllerAdvice
@Slf4j
public class RestExceptionHandler extends ResponseEntityExceptionHandler {
    @Autowired
    ErrorRepository errorRepository;

    @ExceptionHandler(ExpectedException.class)
    // when exception is expected, we should use its HttpStatus when returning ErrorResponse
    private ResponseEntity<ErrorResponse> buildResponseEntityForExpectedException(ExpectedException exception) {
        log.error("Expected exception occurred with code {} : ", exception.getErrorDetails().getCode(), exception.getCause());
        return new ResponseEntity<>(buildErrorResponse(exception, exception.getErrorDetails()), exception.getErrorDetails().getHttpStatus());
    }

    @ExceptionHandler(Exception.class)
    // when exception is not expected, we should use HttpStatus.INTERNAL_SERVER_ERROR when returning ErrorResponse
    private ResponseEntity<ErrorResponse> buildResponseEntityForException(Exception exception) {
        log.error("Unexpected exception occurred : ", exception);
        return new ResponseEntity<>(
                buildErrorResponse(exception, ErrorDetails.ERR_UNEXPECTED_ERROR_OCCURRED.addingArgs(exception.getMessage())),
                HttpStatus.INTERNAL_SERVER_ERROR);
    }

    private ErrorResponse buildErrorResponse(Exception exception, ErrorDetails errorDetails) {
        errorRepository.save(new Error(exception, errorDetails));
        return new ErrorResponse(
                exception.getMessage(),
                errorDetails.getCode(),
                errorDetails.getErrorOnObject() != null ? errorDetails.getErrorOnObject().name() : null,
                errorDetails.getErrorOnField() != null ? errorDetails.getErrorOnField().name() : null,
                exception.getCause() != null ? exception.getCause().getMessage() : null,
                exception.getCause() != null ? exception.getCause().getStackTrace() : null
        );
    }
}
