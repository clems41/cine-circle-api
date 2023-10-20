package com.teasy.CineCircleApi.models.externals;

import info.movito.themoviedbapi.model.*;
import info.movito.themoviedbapi.model.keywords.Keyword;
import info.movito.themoviedbapi.model.people.PersonCast;
import info.movito.themoviedbapi.model.people.PersonCrew;

import java.util.List;

public interface TheMovieDbMedia {
    public int getId();

    public String getBackdropPath();

    public float getPopularity();

    public String getPosterPath();

    public List<Genre> getGenres();

    public String getHomepage();

    public String getOverview();

    public float getVoteAverage();

    public int getVoteCount();

    public String getStatus();

    public List<PersonCast> getCast();

    public List<PersonCrew> getCrew();

    public List<Reviews> getReviews();

    public Credits getCredits();

    public float getUserRating();
}
