package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import jakarta.persistence.Id;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import java.io.Serial;
import java.io.Serializable;
import java.util.Collection;
import java.util.Objects;
import java.util.Set;
import java.util.UUID;

@Entity
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(name = "users",
        indexes = {
                @Index(columnList = "username", unique = true),
                @Index(columnList = "email", unique = true),
                @Index(columnList = "password")
        }
)
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Column(nullable = false, unique = true)
    private String username;

    @Column(nullable = false, unique = true)
    private String email;

    @Column(name = "password", nullable = false)
    private String hashPassword;

    @Column(nullable = false)
    private Boolean enabled;

    @Column(nullable = false)
    private String displayName;

    private String resetPasswordToken;

    @Column(unique = true, updatable = false)
    private UUID topicName;

    @PrePersist
    protected void onCreateAbstractBaseEntity() {
        this.topicName = UUID.randomUUID();
    }

    public User(String username, String email, String hashPassword, String displayName) {
        this.username = username;
        this.email = email;
        this.hashPassword = hashPassword;
        this.displayName = displayName;
        this.enabled = true;
    }
}

