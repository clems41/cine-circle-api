package com.teasy.CineCircleApi.config;

import com.teasy.CineCircleApi.models.entities.Error;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ErrorResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.repositories.ErrorRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

@Order(Ordered.HIGHEST_PRECEDENCE)
@ControllerAdvice
@Slf4j
public class RestExceptionHandler extends ResponseEntityExceptionHandler {
    private final ErrorRepository errorRepository;

    public RestExceptionHandler(ErrorRepository errorRepository) {
        this.errorRepository = errorRepository;
    }

    @ExceptionHandler(ExpectedException.class)
    // when exception is expected, we should use its HttpStatus when returning ErrorResponse
    private ResponseEntity<ErrorResponse> buildResponseEntityForExpectedException(ExpectedException exception) {
        log.error("Expected exception occurred with code {} : ", exception.getErrorDetails().getCode(), exception.getCause());
        return new ResponseEntity<>(buildErrorResponse(exception, exception.getErrorDetails()), exception.getErrorDetails().getHttpStatus());
    }

    // when exception is thrown by jakarta valid annotations
    @Override
    public ResponseEntity<Object> handleMethodArgumentNotValid( MethodArgumentNotValidException exception, HttpHeaders headers, HttpStatusCode status, WebRequest request) {
        try {
            var error = exception.getBindingResult().getFieldErrors().getFirst();
            ErrorDetails errorDetails = ErrorDetails.valueOf(error.getDefaultMessage());
            var expectedException = new ExpectedException(errorDetails, exception);
            return new ResponseEntity<>(
                    buildErrorResponse(expectedException, errorDetails.addingArgs(error.getRejectedValue())),
                    errorDetails.getHttpStatus()
            );
        } catch (IllegalArgumentException e) {
            return new ResponseEntity<>(
                    buildErrorResponse(e, ErrorDetails.ERR_CANNOT_FIND_ERROR_CODE_FROM_VALIDATION_MESSAGE.addingArgs(exception.getMessage())),
                    HttpStatus.INTERNAL_SERVER_ERROR);
        } catch (Exception e) {
            return responseForUnexpectedException(exception);
        }
    }

    // when exception is thrown when converting data
    @ExceptionHandler(MethodArgumentTypeMismatchException.class)
    public ResponseEntity<ErrorResponse> handleMethodArgumentNotValid(MethodArgumentTypeMismatchException exception) {
        var errorDetails = ErrorDetails.ERR_GLOBAL_INVALID_PARAMETER.addingArgs(exception.getName(), exception.getValue());
        var expectedException = new ExpectedException(errorDetails, exception);
        return new ResponseEntity<>(
                buildErrorResponse(expectedException, errorDetails),
                errorDetails.getHttpStatus()
        );
    }

    @ExceptionHandler(Exception.class)
    // when exception is not expected, we should use HttpStatus.INTERNAL_SERVER_ERROR when returning ErrorResponse
    private ResponseEntity<Object> buildResponseEntityForException(Exception exception) {
        return responseForUnexpectedException(exception);
    }

    private ResponseEntity<Object> responseForUnexpectedException(Exception exception) {
        log.error("Unexpected exception occurred : ", exception);
        return new ResponseEntity<>(
                buildErrorResponse(exception, ErrorDetails.ERR_UNEXPECTED_ERROR_OCCURRED.addingArgs(exception.getMessage())),
                HttpStatus.INTERNAL_SERVER_ERROR);
    }

    private ErrorResponse buildErrorResponse(Exception exception, ErrorDetails errorDetails) {
        errorRepository.save(new Error(exception, errorDetails));
        return new ErrorResponse(
                errorDetails.getMessage(),
                errorDetails.getCode(),
                errorDetails.getErrorOnObject() != null ? errorDetails.getErrorOnObject().name() : null,
                errorDetails.getErrorOnField() != null ? errorDetails.getErrorOnField().name() : null,
                exception.getCause() != null ? exception.getCause().getMessage() : null,
                exception.getCause() != null ? exception.getCause().getStackTrace() : null
        );
    }
}
