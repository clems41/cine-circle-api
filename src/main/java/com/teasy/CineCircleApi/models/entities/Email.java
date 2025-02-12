package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.UUID;

@Entity
@Getter
@Setter
@Table(name = "emails",
        indexes = {
                @Index(columnList = "errorCode"),
        }
)
public class Email {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

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

    public Email() {
        this.sent = false;
    }
}

