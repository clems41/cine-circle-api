package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;
import java.util.UUID;

@Getter
@Entity
@Setter
@Table(name = "recommendations",
        indexes = {
                @Index(columnList = "sent_by"),
                @Index(columnList = "receiver"),
                @Index(columnList = "media_id"),
                @Index(columnList = "sentAt DESC")
        }
)
@NoArgsConstructor
public class Recommendation {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Column
    private UUID recommendationRef;

    @ManyToOne
    @JoinColumn(name = "sent_by", referencedColumnName = "id", nullable = false)
    private User sentBy;

    @ManyToOne
    @JoinColumn(name = "media_id", referencedColumnName = "id", nullable = false)
    private Media media;

    @ManyToOne
    @JoinColumn(name = "receiver", referencedColumnName = "id", nullable = false)
    private User receiver;

    @Column
    private String comment;

    @Column
    private Integer rating;

    @Column(nullable = false)
    private Boolean read;

    @Column(nullable = false)
    private LocalDateTime sentAt;

    public Recommendation(
            UUID recommendationRef,
            User sentBy,
            Media media,
            User receiver,
            String comment,
            Integer rating) {
        this.recommendationRef = recommendationRef;
        this.sentBy = sentBy;
        this.media = media;
        this.receiver = receiver;
        this.comment = comment;
        this.rating = rating;
        this.sentAt = LocalDateTime.now();
        this.read = false;
    }
}
