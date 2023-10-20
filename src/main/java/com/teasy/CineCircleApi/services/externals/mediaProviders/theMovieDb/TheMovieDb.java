package com.teasy.CineCircleApi.services.externals.mediaProviders.theMovieDb;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.requests.SearchMediaRequest;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.enums.MediaType;
import com.teasy.CineCircleApi.models.externals.TheMovieDbMedia;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.utils.CustomHttpClient;
import info.movito.themoviedbapi.TmdbApi;
import info.movito.themoviedbapi.TmdbMovies;
import info.movito.themoviedbapi.TmdbTV;
import info.movito.themoviedbapi.model.MovieDb;
import info.movito.themoviedbapi.model.Multi;
import info.movito.themoviedbapi.model.core.NamedIdElement;
import info.movito.themoviedbapi.model.people.PersonCast;
import info.movito.themoviedbapi.model.people.PersonCrew;
import info.movito.themoviedbapi.model.tv.TvSeries;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.Pageable;
import org.springframework.format.datetime.DateFormatter;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.text.ParseException;
import java.util.*;

@Service
@Slf4j
public class TheMovieDb implements MediaProvider {
    private final static String stringArrayDelimiter = ",";
    private final static String language = "fr-FR";
    private final static String jobDirector = "Director";
    private final static String imageUrlPrefix = "https://image.tmdb.org/t/p/w500";

    @Value("${the-movie-db.token}")
    private String token;
    MediaRepository mediaRepository;
    CustomHttpClient customHttpClient;

    @Autowired
    public TheMovieDb(MediaRepository mediaRepository,
                      CustomHttpClient customHttpClient) {
        this.mediaRepository = mediaRepository;
        this.customHttpClient = customHttpClient;
    }

    @Override
    public List<MediaDto> searchMedia(Pageable pageable, SearchMediaRequest searchMediaRequest) {
        var tmdbApi = new TmdbApi(token);
        var multiResponse = tmdbApi.getSearch()
                .searchMulti(searchMediaRequest.query(), language, pageable.getPageNumber())
                .getResults();
        List<MediaDto> result = new ArrayList<>();
        multiResponse.forEach(multi -> {
            Media media = new Media();
            if(multi.getMediaType() == Multi.MediaType.MOVIE) {
                MovieDb movie = (MovieDb) multi;
                media = fromMovieDbToMediaEntity(movie);
            } else if(multi.getMediaType() == Multi.MediaType.TV_SERIES) {
                TvSeries tvSeries = (TvSeries) multi;
                media = fromMTvSeriesToMediaEntity(tvSeries);
            }
            // store result in database with internalId if not already exists
            var existingMedia = findMediaWithExternalId(media.getExternalId());
            if (existingMedia.isEmpty()) {
                var newMedia = mediaRepository.save(media);
                result.add(fromMediaEntityToMediaDto(newMedia));
            } else {
                result.add(fromMediaEntityToMediaDto(existingMedia.get()));
            }

        });
        return result;
    }

    @Override
    public MediaDto getMedia(Long id) throws ResponseStatusException {
        // build example matcher with id
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setId(id);

        // get media from database
        var media = mediaRepository
                .findOne(Example.of(exampleMedia, matcher))
                .orElseThrow(() -> new ResponseStatusException(
                        HttpStatus.NOT_FOUND,
                        String.format("media with id %d cannot be found", id)));
        if (!media.getCompleted()) {
            completeMedia(media);
        }
        return fromMediaEntityToMediaDto(media);
    }

    private Media fromMovieDbToMediaEntity(MovieDb movie) {
        var media = fromMediaInterfaceToMediaEntity((TheMovieDbMedia) movie);
        media.setMediaType(MediaType.MOVIE.name());
        media.setOriginalTitle(movie.getOriginalTitle());
        media.setTitle(movie.getTitle());
        var dateFormatter = new DateFormatter();
        try {
            media.setReleaseDate(dateFormatter.parse(movie.getReleaseDate(), Locale.FRANCE));
        } catch (ParseException e) {
            log.error("cannot parse date {}", movie.getReleaseDate());
        }
        media.setOriginalLanguage(movie.getOriginalLanguage());
        if (media.getCompleted() == null) {
            media.setCompleted(false);
        }
        return media;
    }

    private Media fromMTvSeriesToMediaEntity(TvSeries tvSeries) {
        var media = fromMediaInterfaceToMediaEntity((TheMovieDbMedia) tvSeries);
        media.setMediaType(MediaType.TV_SHOW.name());
        media.setOriginalTitle(tvSeries.getOriginalName());
        media.setTitle(tvSeries.getName());
        var dateFormatter = new DateFormatter();
        try {
            media.setReleaseDate(dateFormatter.parse(tvSeries.getFirstAirDate(), Locale.FRANCE));
        } catch (ParseException e) {
            log.error("cannot parse date {}", tvSeries.getFirstAirDate());
        }
        if (tvSeries.getOriginCountry() != null) {
            media.setOriginCountry(String.join(stringArrayDelimiter,
                    tvSeries.getOriginCountry()
                            .stream()
                            .filter(s -> !s.isEmpty())
                            .toList()
            ));
        }
        if (media.getCompleted() == null) {
            media.setCompleted(false);
        }
        return media;
    }

    private Media fromMediaInterfaceToMediaEntity(TheMovieDbMedia mediaInterface) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        Media media = mapper.convertValue(mediaInterface, Media.class);
        media.setExternalId(String.valueOf(mediaInterface.getId()));
        media.setMediaProvider(com.teasy.CineCircleApi.models.enums.MediaProvider.THE_MOVIE_DATABASE.name());
        media.setPosterUrl(getCompleteImageUrl(mediaInterface.getPosterPath()));
        media.setBackdropUrl(getCompleteImageUrl(mediaInterface.getBackdropPath()));
        media.setVoteAverage(mediaInterface.getVoteAverage());
        media.setVoteCount(mediaInterface.getVoteCount());
        if (mediaInterface.getGenres() != null) {
            media.setGenres(String.join(stringArrayDelimiter,
                    mediaInterface.getGenres()
                            .stream()
                            .map(NamedIdElement::getName)
                            .filter(s -> !s.isEmpty())
                            .toList()
            ));
        }
        if (media.getCompleted() == null) {
            media.setCompleted(false);
        }
        return media;
    }

    private void completeMedia(Media media) {
        var tmdbApi = new TmdbApi(token);
        // get casting
        List<PersonCast> cast = new ArrayList<>();
        List<PersonCrew> crew = new ArrayList<>();
        if (Objects.equals(media.getMediaType(), MediaType.MOVIE.name())) {
            MovieDb movie = tmdbApi.getMovies()
                    .getMovie(Integer.parseInt(media.getExternalId()), language, TmdbMovies.MovieMethod.credits);
            cast = movie.getCredits().getCast();
            crew = movie.getCredits().getCrew();
        }
        if (Objects.equals(media.getMediaType(), MediaType.TV_SHOW.name())) {
            TvSeries tvSeries = tmdbApi.getTvSeries()
                    .getSeries(Integer.parseInt(media.getExternalId()), language, TmdbTV.TvMethod.credits);
            cast = tvSeries.getCredits().getCast();
            crew = tvSeries.getCredits().getCrew();
        }

        // adding cast and crew to media
        if (cast != null && !cast.isEmpty()) {
            media.setDirector(getOnlyDirectorsFromCrew(crew));
            media.setActors(getOnlyActorsFromCast(cast));
        }

        // mark media as completed to avoid getting details again later
        media.setCompleted(true);
        mediaRepository.save(media);
    }

    private String getOnlyActorsFromCast(List<PersonCast> cast) {
        return String.join(stringArrayDelimiter, cast
                .stream()
                .map(PersonCast::getName)
                .filter(s -> !s.isEmpty())
                .toList()
        );
    }

    private String getOnlyDirectorsFromCrew(List<PersonCrew> cast) {
        return String.join(stringArrayDelimiter, cast
                .stream()
                .filter(personCrew -> jobDirector.equals(personCrew.getJob())) // only actors
                .map(PersonCrew::getName)
                .filter(s -> !s.isEmpty())
                .toList()
        );
    }

    private Optional<Media> findMediaWithExternalId(String externalId) {
        // build example matcher with external id and media provider
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setExternalId(String.valueOf(externalId));
        exampleMedia.setMediaProvider(com.teasy.CineCircleApi.models.enums.MediaProvider.THE_MOVIE_DATABASE.name());

        return mediaRepository
                .findOne(Example.of(exampleMedia, matcher));
    }

    private String getCompleteImageUrl(String posterUrl) {
        return String.format("%s%s", imageUrlPrefix, posterUrl);
    }

    private MediaDto fromMediaEntityToMediaDto(Media media) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        return mapper.convertValue(media, MediaDto.class);
    }
}
