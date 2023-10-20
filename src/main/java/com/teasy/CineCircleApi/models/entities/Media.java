package com.teasy.CineCircleApi.models.entities;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.teasy.CineCircleApi.models.enums.MediaProvider;
import com.teasy.CineCircleApi.models.enums.MediaType;
import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

@Entity
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(name = "medias")
public class Media {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String externalId;

    @Column(nullable = false)
    private String mediaProvider;

    private String title;

    private String originalTitle;

    @Column(columnDefinition = "TEXT")
    private String overview;

    private String posterUrl;

    private String backdropUrl;

    @Column(columnDefinition = "TEXT")
    private String genres;

    private String mediaType;

    private Date releaseDate;

    private String originalLanguage;

    private Float popularity;

    private Float voteAverage;

    private Integer voteCount;

    private String originCountry;

    private Boolean completed;

    @Column(columnDefinition = "TEXT")
    private String actors;

    @Column(columnDefinition = "TEXT")
    private String director;
}
