package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.annotation.LastModifiedDate;

import java.time.Instant;
import java.util.UUID;

@Getter
@Entity
@Setter
@Table(name = "headings",
        indexes = {
                @Index(columnList = "user_id"),
                @Index(columnList = "media_id"),
                @Index(columnList = "added_at DESC"),
                @Index(columnList = "user_id, added_at")
        }
)
@NoArgsConstructor
public class Heading {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    UUID id;

    @ManyToOne
    @JoinColumn(name = "user_id")
    User user;

    @ManyToOne
    @JoinColumn(name = "media_id")
    Media media;

    @Column(name = "added_at")
    @LastModifiedDate
    Instant addedAt;

    public Heading(User user, Media media) {
        this.user = user;
        this.media = media;
    }
}
