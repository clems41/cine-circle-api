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
@Table(name = "recommendations")
@NoArgsConstructor
public class Recommendation {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @ManyToOne(cascade = {CascadeType.PERSIST, CascadeType.MERGE, CascadeType.REFRESH})
    @JoinColumn(name = "sent_by", referencedColumnName = "id", nullable=false)
    private User sentBy;

    @ManyToOne(cascade = {CascadeType.PERSIST, CascadeType.MERGE, CascadeType.REFRESH})
    @JoinColumn(name = "media_id", referencedColumnName = "id", nullable=false)
    private Media media;

    @ManyToMany
    @JoinTable(
            name = "recommendation_users",
            joinColumns = @JoinColumn(name = "recommendation_id"),
            inverseJoinColumns = @JoinColumn(name = "user_id"))
    private Set<User> receivers;

    @Column(name = "comment", nullable = false)
    private String comment;

    @Column(name = "rating", nullable = false)
    private Integer rating;

    @Column(name = "sent_at", nullable = false)
    private LocalDateTime sentAt;

    public Recommendation(
            User sentBy,
            Media media,
            Set<User> receivers,
            String comment,
            Integer rating) {
        this.sentBy = sentBy;
        this.media = media;
        this.receivers = receivers;
        this.comment = comment;
        this.rating = rating;
        this.sentAt = LocalDateTime.now();
    }
}
