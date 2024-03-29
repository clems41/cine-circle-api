package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;
import java.util.Set;
import java.util.UUID;

@Getter
@Entity
@Setter
@Table(name = "recommendations",
        indexes = {
                @Index(columnList = "sent_by"),
                @Index(columnList = "media_id"),
                @Index(columnList = "sentAt DESC")
        }
)
@NoArgsConstructor
public class Recommendation {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @ManyToOne(cascade = {CascadeType.PERSIST, CascadeType.MERGE, CascadeType.REFRESH})
    @JoinColumn(name = "sent_by", referencedColumnName = "id", nullable = false)
    private User sentBy;

    @ManyToOne(cascade = {CascadeType.PERSIST, CascadeType.MERGE, CascadeType.REFRESH})
    @JoinColumn(name = "media_id", referencedColumnName = "id", nullable = false)
    private Media media;

    @ManyToMany
    @JoinTable(
            name = "recommendation_users",
            joinColumns = @JoinColumn(name = "recommendation_id"),
            inverseJoinColumns = @JoinColumn(name = "user_id"))
    private Set<User> receivers;

    @ManyToMany
    @JoinTable(
            name = "recommendation_circles",
            joinColumns = @JoinColumn(name = "recommendation_id"),
            inverseJoinColumns = @JoinColumn(name = "circle_id"))
    private Set<Circle> circles;

    @Column(nullable = false)
    private String comment;

    @Column(nullable = false)
    private Integer rating;

    @Column(nullable = false)
    private LocalDateTime sentAt;

    public Recommendation(
            User sentBy,
            Media media,
            Set<User> receivers,
            Set<Circle> circles,
            String comment,
            Integer rating) {
        this.sentBy = sentBy;
        this.media = media;
        this.receivers = receivers;
        this.circles = circles;
        this.comment = comment;
        this.rating = rating;
        this.sentAt = LocalDateTime.now();
    }
}
