package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.relational.core.sql.In;

import java.time.LocalDateTime;
import java.util.UUID;

@Getter
@Entity
@Setter
@Table(name = "libraries")
@NoArgsConstructor
public class Library {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    UUID id;

    @ManyToOne
    @JoinColumn(name = "user_id")
    User user;

    @ManyToOne
    @JoinColumn(name = "media_id")
    Media media;


    @Column(name = "added_at", nullable = false)
    LocalDateTime addedAt;

    String comment;

    Integer rating;

    public Library(User user, Media media, String comment, Integer rating) {
        this.user = user;
        this.media = media;
        this.comment = comment;
        this.rating = rating;
        this.addedAt = LocalDateTime.now();
    }
}
