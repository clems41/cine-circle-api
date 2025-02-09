package com.teasy.CineCircleApi.models.exceptions;


import lombok.Getter;

@Getter
public class ExpectedException extends Exception {
    private final ErrorDetails errorDetails;

    public ExpectedException(ErrorDetails error) {
        super(error.getMessage());
        this.errorDetails = error;
    }

    public ExpectedException(ErrorDetails error, Throwable cause) {
        super(error.getMessage(), cause);
        this.errorDetails = error;
    }
}
