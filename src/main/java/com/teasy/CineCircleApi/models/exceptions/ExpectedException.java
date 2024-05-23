package com.teasy.CineCircleApi.models.exceptions;


import com.teasy.CineCircleApi.models.enums.ErrorMessage;
import lombok.Getter;
import org.springframework.http.HttpStatus;

@Getter
public class ExpectedException extends Exception {
    private final HttpStatus httpStatus;
    private final ErrorMessage error;

    public ExpectedException(ErrorMessage error, HttpStatus httpStatus) {
        super(error.getMessage());
        this.httpStatus = httpStatus;
        this.error = error;
    }

    public ExpectedException(ErrorMessage error, Throwable cause, HttpStatus httpStatus) {
        super(error.getMessage(), cause);
        this.httpStatus = httpStatus;
        this.error = error;
    }

    public ExpectedException(ErrorMessage errorForExpectedException) {
        this(errorForExpectedException, HttpStatus.INTERNAL_SERVER_ERROR);
    }

    public ExpectedException(ErrorMessage errorForExpectedException, Throwable cause) {
        this(errorForExpectedException, cause, HttpStatus.INTERNAL_SERVER_ERROR);
    }
}
