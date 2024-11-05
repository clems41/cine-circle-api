package com.teasy.CineCircleApi.services.externals.mediaProviders.theMovieDb;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.enums.ErrorMessage;
import com.teasy.CineCircleApi.models.enums.MediaProviderEnum;
import com.teasy.CineCircleApi.models.enums.MediaTypeEnum;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.externals.ExternalMedia;
import com.teasy.CineCircleApi.models.externals.ExternalMediaShort;
import com.teasy.CineCircleApi.models.externals.theMovieDb.WatchProviderInfo;
import com.teasy.CineCircleApi.models.externals.theMovieDb.WatchProvidersResponse;
import com.teasy.CineCircleApi.models.utils.CustomHttpClientSendRequest;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.services.utils.CustomHttpClient;
import info.movito.themoviedbapi.TmdbApi;
import info.movito.themoviedbapi.TmdbMovies;
import info.movito.themoviedbapi.TmdbTV;
import info.movito.themoviedbapi.model.Genre;
import info.movito.themoviedbapi.model.MovieDb;
import info.movito.themoviedbapi.model.Multi;
import info.movito.themoviedbapi.model.Video;
import info.movito.themoviedbapi.model.core.NamedIdElement;
import info.movito.themoviedbapi.model.people.Person;
import info.movito.themoviedbapi.model.people.PersonCast;
import info.movito.themoviedbapi.model.people.PersonCrew;
import info.movito.themoviedbapi.model.tv.TvSeries;
import info.movito.themoviedbapi.tools.ApiUrl;
import info.movito.themoviedbapi.tools.RequestMethod;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.*;

@Service
@Slf4j
public class TheMovieDbService implements MediaProvider {
    private final static String stringArrayDelimiter = ",";

    private final static String theMovieDbApiBaseUrl = "https://api.themoviedb.org/3";
    private final static String tvSuffix = "/tv";
    private final static String movieSuffix = "/movie";
    private final static String watchProvidersSuffix = "/watch/providers";
    private final static String language = "fr-FR";
    private final static String jobDirector = "Director";
    private final static String imageUrlPrefix = "https://image.tmdb.org/t/p/w500";

    private final static String youtubeVideoUrlPrefix = "https://www.youtube.com/watch?v=";

    private final static String youtubeVideoSite = "YouTube";
    private final static String trailerVideoType = "Trailer";

    private final TmdbApi tmdbApi;

    public TheMovieDbService(@Value("${the-movie-db.api-key}") String apiKey) {
        tmdbApi = new TmdbApi(apiKey);
    }

    @Override
    public List<ExternalMediaShort> searchMedia(Pageable pageable,
                                                MediaSearchRequest mediaSearchRequest) {
        var multiResponse = tmdbApi.getSearch()
                .searchMulti(mediaSearchRequest.query(), language, pageable.getPageNumber())
                .getResults();
        List<ExternalMediaShort> result = new ArrayList<>();

        multiResponse.forEach(multi -> {
            ExternalMediaShort media;
            if (multi.getMediaType() == Multi.MediaType.MOVIE) {
                MovieDb movie = (MovieDb) multi;
                media = fromMovieDbMediaToExternalMediaShort(movie, MovieDb.class);
            } else if (multi.getMediaType() == Multi.MediaType.TV_SERIES) {
                TvSeries tvSeries = (TvSeries) multi;
                media = fromMovieDbMediaToExternalMediaShort(tvSeries, TvSeries.class);
            } else {
                return;
            }
            result.add(media);
        });
        return result;
    }

    @Override
    public ExternalMedia getMedia(String externalId, MediaTypeEnum mediaType) throws ExpectedException {
        // get casting
        List<PersonCast> cast;
        List<PersonCrew> crew;
        List<Person> persons;
        List<Genre> genres;
        List<Video> videos;
        Integer runtime = null;
        if (Objects.equals(mediaType, MediaTypeEnum.MOVIE)) {
            MovieDb movie = tmdbApi.getMovies()
                    .getMovie(Integer.parseInt(externalId),
                            language,
                            TmdbMovies.MovieMethod.credits,
                            TmdbMovies.MovieMethod.videos);
            cast = movie.getCredits().getCast();
            crew = movie.getCredits().getCrew();
            genres = movie.getGenres();
            videos = movie.getVideos();
            persons = new ArrayList<>();
            runtime = movie.getRuntime();
        } else if (Objects.equals(mediaType, MediaTypeEnum.TV_SHOW)) {
            TvSeries tvSeries = tmdbApi.getTvSeries()
                    .getSeries(Integer.parseInt(externalId),
                            language,
                            TmdbTV.TvMethod.credits,
                            TmdbTV.TvMethod.videos);
            cast = tvSeries.getCredits().getCast();
            persons = tvSeries.getCreatedBy();
            genres = tvSeries.getGenres();
            videos = tvSeries.getVideos();
            crew = new ArrayList<>();
            runtime = tvSeries.getEpisodeRuntime() != null && !tvSeries.getEpisodeRuntime().isEmpty() ?
                    tvSeries.getEpisodeRuntime().get(0) : null;
        } else {
            throw new ExpectedException(
                    ErrorMessage.MEDIA_NOT_FOUND,
                    HttpStatus.NOT_FOUND
            );
        }
        var media = new ExternalMedia();

        // adding cast and crew to media
        if (cast != null && !cast.isEmpty()) {
            media.setActors(getOnlyActorsFromCast(cast));
        }
        if (crew != null && !crew.isEmpty()) {
            media.setDirector(getOnlyDirectorsFromCrew(crew));
        }
        if (persons != null && !persons.isEmpty()) {
            media.setDirector(getOnlyDirectorsFromPersons(persons));
        }

        // adding genre
        if (genres != null && !genres.isEmpty()) {
            media.setGenres(String.join(stringArrayDelimiter, genres
                    .stream()
                    .map(NamedIdElement::getName)
                    .filter(s -> !s.isEmpty())
                    .toList()
            ));
        }

        // adding trailer
        if (videos != null && !videos.isEmpty()) {
            media.setTrailerUrl(getTrailerUrl(videos));
        }

        // adding runtime
        media.setRuntime(runtime);
        return media;
    }

    @Override
    public List<String> listGenres() {
        var genres = tmdbApi.getGenre().getGenreList(language);
        return genres.stream().map(NamedIdElement::getName).toList();
    }

    @Override
    public MediaProviderEnum getMediaProvider() {
        return MediaProviderEnum.THE_MOVIE_DATABASE;
    }

    @Override
    public List<String> getWatchProvidersForMedia(String externalId, MediaTypeEnum mediaType) throws ExpectedException {
        // define url depending on media type and id
        String url = "";
        if (mediaType.equals(MediaTypeEnum.MOVIE)) {
            url = url.concat(movieSuffix);
        } else {
            url = url.concat(tvSuffix);
        }
        url = url.concat("/").concat(externalId).concat(watchProvidersSuffix).substring(1);
        ApiUrl apiUrl = new ApiUrl(url);
        var response = tmdbApi.requestWebPage(apiUrl, null, RequestMethod.GET);
        ObjectMapper mapper = new ObjectMapper()
                .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        try {
            var watchProvidersResponse = mapper.readValue(response, WatchProvidersResponse.class);
            if (watchProvidersResponse.getResults() == null) {
                return new ArrayList<>();
            }
            return watchProvidersResponse.getResults().getFr().getFlatrate().stream().map(WatchProviderInfo::getProviderName).toList();
        } catch (Exception e) {
            log.warn("Error when getting watch providers with url {} : {}", url, e.getMessage());
            throw new ExpectedException(
                    ErrorMessage.MEDIA_NOT_FOUND,
                    HttpStatus.NOT_FOUND
            );
        }
    }

    private <T> ExternalMediaShort fromMovieDbMediaToExternalMediaShort(T movieDbMedia, Class<T> movieDbMediaType) {
        if (movieDbMediaType != MovieDb.class && movieDbMediaType != TvSeries.class) {
            return null;
        }
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        ExternalMediaShort media = mapper.convertValue(movieDbMedia, ExternalMediaShort.class);
        if (movieDbMediaType == MovieDb.class) {
            MovieDb movie = (MovieDb) movieDbMedia;
            media.setExternalId(String.valueOf(movie.getId()));
            media.setPosterUrl(getCompleteImageUrl(movie.getPosterPath()));
            media.setBackdropUrl(getCompleteImageUrl(movie.getBackdropPath()));
            media.setMediaType(MediaTypeEnum.MOVIE.name());
            media.setOriginalTitle(movie.getOriginalTitle());
            media.setTitle(movie.getTitle());
            if (movie.getReleaseDate() != null && !movie.getReleaseDate().isEmpty()) {
                media.setReleaseDate(LocalDate.parse(movie.getReleaseDate()));
            }
            media.setOriginalLanguage(movie.getOriginalLanguage());
        } else {
            TvSeries tvSeries = (TvSeries) movieDbMedia;
            media.setExternalId(String.valueOf(tvSeries.getId()));
            media.setPosterUrl(getCompleteImageUrl(tvSeries.getPosterPath()));
            media.setBackdropUrl(getCompleteImageUrl(tvSeries.getBackdropPath()));
            media.setMediaType(MediaTypeEnum.TV_SHOW.name());
            media.setOriginalTitle(tvSeries.getOriginalName());
            media.setTitle(tvSeries.getName());
            if (tvSeries.getFirstAirDate() != null && !tvSeries.getFirstAirDate().isEmpty()) {
                media.setReleaseDate(LocalDate.parse(tvSeries.getFirstAirDate()));
            }
        }
        return media;
    }

    private String getTrailerUrl(List<Video> videos) {
        String youtubeVideoKey = videos
                .stream()
                .filter(video -> Objects.equals(video.getSite(), youtubeVideoSite) &&
                        Objects.equals(video.getType(), trailerVideoType))
                .map(Video::getKey)
                .toList()
                .get(0);
        return youtubeVideoUrlPrefix.concat(youtubeVideoKey);
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

    private String getOnlyDirectorsFromPersons(List<Person> persons) {
        return String.join(stringArrayDelimiter, persons
                .stream()
                .map(Person::getName)
                .filter(s -> !s.isEmpty())
                .toList()
        );
    }

    private String getCompleteImageUrl(String posterUrl) {
        if (posterUrl == null || posterUrl.isEmpty()) {
            return null;
        }
        return String.format("%s%s", imageUrlPrefix, posterUrl);
    }
}
