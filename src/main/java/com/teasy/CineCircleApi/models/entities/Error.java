package com.teasy.CineCircleApi.models.entities;

import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.utils.StringUtils;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Index;
import jakarta.persistence.Table;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;
import java.util.Arrays;

@Entity
@Getter
@Setter
@Table(name = "errors",
        indexes = {
                @Index(columnList = "code"),
        }
)
@NoArgsConstructor
public class Error extends BaseEntity {
    @Column
    private Integer httpStatusCode;

    @Column
    private String message;

    @Column
    private String code;

    @Column
    private String object;

    @Column
    private String field;

    @Column
    private String cause;

    @Column
    private String firstElementOfStackTrace;

    @Column(nullable = false)
    private LocalDateTime triggeredAt;

    public Error(Exception exception, ErrorDetails errorDetails) {
        this.triggeredAt = LocalDateTime.now();
        if (errorDetails != null) {
            this.httpStatusCode = errorDetails.getHttpStatus().value();
            this.message = StringUtils.substringForDatabase(errorDetails.getMessage());
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
                this.cause = StringUtils.substringForDatabase(exception.getCause().getMessage());
            }
            if (exception.getCause().getStackTrace() != null) {
                var stackTrace = Arrays.stream(exception.getCause().getStackTrace()).toList();
                if (!stackTrace.isEmpty()) {
                    var firstElementOfStackTrace = stackTrace.getFirst().toString();
                    this.firstElementOfStackTrace = StringUtils.substringForDatabase(firstElementOfStackTrace);
                }
            }
        }
    }
}

