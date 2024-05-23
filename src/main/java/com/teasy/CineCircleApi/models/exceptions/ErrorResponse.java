package com.teasy.CineCircleApi.models.exceptions;

public record ErrorResponse(
    String errorMessage,
    String errorCode,
    String errorCause,
    StackTraceElement[] errorStack
){}
