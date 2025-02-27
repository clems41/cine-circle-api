package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;

@Getter
@Entity
@Setter
@Table(name = "libraries",
        indexes = {
                @Index(columnList = "user_id"),
                @Index(columnList = "media_id"),
                @Index(columnList = "addedAt DESC"),
                @Index(columnList = "user_id, addedAt")
        }
)
@NoArgsConstructor
public class Library extends BaseEntity {
    @ManyToOne
    @JoinColumn(name = "user_id")
    User user;

    @ManyToOne
    @JoinColumn(name = "media_id")
    Media media;


    @Column(nullable = false)
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
