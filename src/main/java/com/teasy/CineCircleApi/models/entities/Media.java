package com.teasy.CineCircleApi.models.entities;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Index;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDate;

@Entity
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Table(name = "medias",
        indexes = {
                @Index(columnList = "externalId, mediaProvider")
        }
)
public class Media extends BaseEntity {
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

    private String trailerUrl;

    @Column(columnDefinition = "TEXT")
    private String genres;

    private MediaTypeEnum mediaType;

    @JsonFormat(pattern = "dd/MM/yyyy")
    private LocalDate releaseDate;

    private Integer runtime;

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
