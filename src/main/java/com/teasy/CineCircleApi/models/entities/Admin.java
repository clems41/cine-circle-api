package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Entity
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(name = "admins",
        indexes = {
                @Index(columnList = "email", unique = true),
        }
)
public class Admin {
    @Id
    @Column(nullable = false, unique = true)
    private String email;

    @Column(nullable = false)
    private Boolean shouldReceiveFeedback;
}

