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

@Entity
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "username", nullable = false, unique = true)
    private String username;

    @Column(name = "email", nullable = false, unique = true)
    private String email;
    
    @Column(name = "password", nullable = false)
    private String hashPassword;

    @Column(name = "enabled", nullable = false)
    private Boolean enabled;

    @Column(name = "display_name")
    private String displayName;

    public User(String username, String email, String hashPassword, String displayName) {
        this.username = username;
        this.email = email;
        this.hashPassword = hashPassword;
        this.displayName = displayName;
        this.enabled = true;
    }
}

