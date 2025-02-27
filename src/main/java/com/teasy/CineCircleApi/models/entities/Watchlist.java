package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;

@Getter
@Entity
@Setter
@Table(name = "watchlists",
        indexes = {
                @Index(columnList = "user_id"),
                @Index(columnList = "media_id"),
                @Index(columnList = "addedAt DESC"),
                @Index(columnList = "user_id, addedAt")
        }
)
@NoArgsConstructor
public class Watchlist extends BaseEntity {
    @ManyToOne
    @JoinColumn(name = "user_id")
    User user;

    @ManyToOne
    @JoinColumn(name = "media_id")
    Media media;


    @Column(nullable = false)
    LocalDateTime addedAt;

    public Watchlist(User user, Media media) {
        this.user = user;
        this.media = media;
        this.addedAt = LocalDateTime.now();
    }
}
