package com.teasy.CineCircleApi.models.entities;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;
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
public class User extends BaseEntity {
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

    @ManyToMany
    @JoinTable(
            name = "user_relationships",
            joinColumns = @JoinColumn(name = "user_id"),
            inverseJoinColumns = @JoinColumn(name = "related_to_user_id")
    )
    @JsonIgnoreProperties(value = "relatedUsers")
    private Set<User> relatedUsers;

    @ManyToMany
    @JoinTable(
            name = "headings",
            joinColumns = @JoinColumn(name = "user_id"),
            inverseJoinColumns = @JoinColumn(name = "media_id"))
    private Set<Media> headings;

    @Column
    private String refreshToken;

    @Column
    private LocalDateTime refreshTokenExpirationDate;

    @PrePersist
    protected void onCreateAbstractBaseEntity() {
        this.topicName = UUID.randomUUID();
    }

    public User(String username, String email, String hashPassword, String displayName, String refreshToken, LocalDateTime refreshTokenExpirationDate) {
        this.username = username;
        this.email = email;
        this.hashPassword = hashPassword;
        this.displayName = displayName;
        this.enabled = true;
        this.relatedUsers = Set.of();
        this.refreshToken = refreshToken;
        this.refreshTokenExpirationDate = refreshTokenExpirationDate;
    }

    public void addRelatedUser(User user) {
        this.relatedUsers.add(user);
    }

    public void removeRelatedUser(User user) {
        this.relatedUsers.remove(user);
    }

    public void addMediaToHeadings(Media media) {
        this.headings.add(media);
    }

    public void removeMediaFromHeadings(Media media) {
        this.headings.remove(media);
    }
}

