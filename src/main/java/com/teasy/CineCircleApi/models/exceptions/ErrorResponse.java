package com.teasy.CineCircleApi.models.exceptions;

public record ErrorResponse(
    String errorMessage,
    String errorCode,
    String errorOnObject,
    String errorOnField,
    String errorCause,
    StackTraceElement[] errorStack
){}
