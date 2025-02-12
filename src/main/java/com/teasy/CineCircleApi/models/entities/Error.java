package com.teasy.CineCircleApi.models.entities;

import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Arrays;
import java.util.UUID;

@Entity
@Getter
@Setter
@Table(name = "errors",
        indexes = {
                @Index(columnList = "errorCode"),
        }
)
@NoArgsConstructor
public class Error {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    UUID id;

    @Column
    Integer httpStatusCode;

    @Column
    String message;

    @Column
    String code;

    @Column
    String object;

    @Column
    String field;

    @Column
    String cause;

    @Column
    String firstElementOfStackTrace;

    public Error(Exception exception, ErrorDetails errorDetails) {
        if (errorDetails != null) {
            this.httpStatusCode = errorDetails.getHttpStatus().value();
            this.message = errorDetails.getMessage();
            this.code = errorDetails.getCode();
            if (errorDetails.getErrorOnObject() != null) {
                this.object = errorDetails.getErrorOnObject().name();
            }
            if (errorDetails.getErrorOnField() != null) {
                this.field = errorDetails.getErrorOnField().name();
            }
        }
        if (exception != null && exception.getCause() != null) {
            if (exception.getCause() != null) {
                this.cause = exception.getCause().getMessage();
            }
            if (exception.getCause().getStackTrace() != null) {
                var stackTrace = Arrays.stream(exception.getCause().getStackTrace()).toList();
                if (!stackTrace.isEmpty()) {
                    this.firstElementOfStackTrace = stackTrace.getFirst().toString();
                }
            }
        }
    }
}

