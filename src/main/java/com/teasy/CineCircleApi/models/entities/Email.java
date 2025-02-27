package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Index;
import jakarta.persistence.Table;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;

@Entity
@Getter
@Setter
@Table(name = "emails",
        indexes = {
                @Index(columnList = "receiver"),
        }
)
public class Email extends BaseEntity {
    @Column
    private String sender;

    @Column
    private String receiver;

    @Column
    private String subject;

    @Column
    private String templateName;

    @Column
    private String templateValues;

    @Column
    private Boolean sent;

    @Column
    private String error;

    @Column(nullable = false)
    private LocalDateTime sentAt;

    public Email() {
        this.sent = false;
        this.sentAt = LocalDateTime.now();
    }
}

