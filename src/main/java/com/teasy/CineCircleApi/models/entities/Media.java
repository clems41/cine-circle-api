package com.teasy.CineCircleApi.models.entities;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Date;

@Entity
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(name = "medias")
public class Media {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long internalId;

    @Column(nullable = false)
    private String externalId;

    @Column(nullable = false)
    private MediaProvider mediaProvider;

    private String title;

    private String originalTitle;

    @Column(columnDefinition = "TEXT")
    private String synopsis;

    private String posterUrl;

    private MediaType mediaType;

    private Date releaseDate;
}
