package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;
import java.util.Date;

@Getter
@Entity
@Setter
@Table(name = "watchlists")
@NoArgsConstructor
public class Watchlist {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    Long id;

    @ManyToOne
    @JoinColumn(name = "user_id")
    User user;

    @ManyToOne
    @JoinColumn(name = "media_id")
    Media media;


    @Column(name = "added_at", nullable = false)
    LocalDateTime addedAt;

    public Watchlist(User user, Media media) {
        this.user = user;
        this.media = media;
        this.addedAt = LocalDateTime.now();
    }
}
