package com.teasy.CineCircleApi.services.externals.mediaProviders.theMovieDb;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.teasy.CineCircleApi.models.dtos.MediaCompleteDto;
import com.teasy.CineCircleApi.models.dtos.MediaDto;
import com.teasy.CineCircleApi.models.dtos.MediaRecommendationDto;
import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.requests.MediaSearchRequest;
import com.teasy.CineCircleApi.models.dtos.requests.RecommendationReceivedRequest;
import com.teasy.CineCircleApi.models.dtos.responses.MediaGenreResponse;
import com.teasy.CineCircleApi.models.entities.Media;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.enums.MediaType;
import com.teasy.CineCircleApi.models.exceptions.CustomException;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.repositories.MediaRepository;
import com.teasy.CineCircleApi.repositories.RecommendationRepository;
import com.teasy.CineCircleApi.repositories.UserRepository;
import com.teasy.CineCircleApi.services.RecommendationService;
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
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.domain.Example;
import org.springframework.data.domain.ExampleMatcher;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.*;

@Service
@Slf4j
public class TheMovieDbService implements MediaProvider {
    private final static String stringArrayDelimiter = ",";
    private final static String language = "fr-FR";
    private final static String jobDirector = "Director";
    private final static String imageUrlPrefix = "https://image.tmdb.org/t/p/w500";

    private final static String youtubeVideoUrlPrefix = "https://www.youtube.com/watch?v=";

    private final static String youtubeVideoSite = "YouTube";
    private final static String trailerVideoType = "Trailer";

    @Value("${the-movie-db.api-key}")
    private String apiKey;

    private TmdbApi tmdbApi;
    MediaRepository mediaRepository;
    RecommendationService recommendationService;
    CustomHttpClient customHttpClient;

    @Autowired
    public TheMovieDbService(MediaRepository mediaRepository,
                             CustomHttpClient customHttpClient,
                             RecommendationService recommendationService) {
        this.mediaRepository = mediaRepository;
        this.customHttpClient = customHttpClient;
        this.recommendationService = recommendationService;
    }

    @Override
    public List<MediaDto> searchMedia(Pageable pageable,
                                      MediaSearchRequest mediaSearchRequest,
                                      String authenticatedUsername) {
        initTmdbApi();
        var multiResponse = tmdbApi.getSearch()
                .searchMulti(mediaSearchRequest.query(), language, pageable.getPageNumber())
                .getResults();
        List<MediaDto> result = new ArrayList<>();
        multiResponse.forEach(multi -> {
            Media media;
            if (multi.getMediaType() == Multi.MediaType.MOVIE) {
                MovieDb movie = (MovieDb) multi;
                media = fromMovieDbToMediaEntity(movie);
            } else if (multi.getMediaType() == Multi.MediaType.TV_SERIES) {
                TvSeries tvSeries = (TvSeries) multi;
                media = fromMTvSeriesToMediaEntity(tvSeries);
            } else {
                return;
            }

            // store result in database with internalId if not already exists
            var existingMedia = findMediaWithExternalId(media.getExternalId());
            if (existingMedia.isEmpty()) {
                var newMedia = mediaRepository.save(media);
                result.add(fromMediaEntityToDto(newMedia, MediaDto.class, authenticatedUsername));
            } else {
                result.add(fromMediaEntityToDto(existingMedia.get(), MediaDto.class, authenticatedUsername));
            }

        });
        return result;
    }

    @Override
    public MediaCompleteDto getMedia(Long id, String authenticatedUsername) throws CustomException {
        // build example matcher with id
        ExampleMatcher matcher = ExampleMatcher
                .matchingAll()
                .withIgnoreNullValues();
        var exampleMedia = new Media();
        exampleMedia.setId(id);

        // get media from database
        var media = mediaRepository
                .findOne(Example.of(exampleMedia, matcher))
                .orElseThrow(() -> CustomExceptionHandler.mediaWithIdNotFound(id));
        if (!media.getCompleted()) {
            completeMedia(media);
        }
        return fromMediaEntityToDto(media, MediaCompleteDto.class, authenticatedUsername);
    }

    @Override
    public MediaGenreResponse listGenres() {
        initTmdbApi();
        var genres = tmdbApi.getGenre().getGenreList(language);
        return new MediaGenreResponse(genres.stream().map(NamedIdElement::getName).toList());
    }

    private void initTmdbApi() {
        if (tmdbApi == null) {
            tmdbApi = new TmdbApi(apiKey);
        }
    }

    private <T> void addRecommendationRatingFields(T mediaDto, String authenticatedUsername) {
        if (mediaDto.getClass() != MediaDto.class && mediaDto.getClass() != MediaCompleteDto.class) {
            return;
        }

        // find recommendation average and count
        Long mediaId;
        if (mediaDto.getClass() == MediaDto.class) {
            mediaId = ((MediaDto) mediaDto).getId();
        } else {
            mediaId = ((MediaCompleteDto) mediaDto).getId();
        }
        var recommendations = findRecommendationsForMediaAndAuthenticatedUsername(mediaId, authenticatedUsername);
        var recommendationRatingCount = recommendations.size();
        var recommendationRatingAverage = recommendations
                .stream()
                .mapToDouble(MediaRecommendationDto::getRating)
                .average();
        if (mediaDto.getClass() == MediaDto.class) {
            ((MediaDto) mediaDto).setRecommendationRatingCount(recommendationRatingCount);
            ((MediaDto) mediaDto).setRecommendationRatingAverage(recommendationRatingAverage.isPresent() ?
                    recommendationRatingAverage.getAsDouble() : null);
        } else {
            ((MediaCompleteDto) mediaDto).setRecommendationRatingCount(recommendationRatingCount);
            ((MediaCompleteDto) mediaDto).setRecommendationRatingAverage(recommendationRatingAverage.isPresent() ?
                    recommendationRatingAverage.getAsDouble() : null);
            // add all recommendations comment for complete media dto
            ((MediaCompleteDto) mediaDto).setRecommendations(recommendations);
        }
    }

    private List<MediaRecommendationDto> findRecommendationsForMediaAndAuthenticatedUsername(Long mediaId, String username) {
        var request = new RecommendationReceivedRequest(mediaId);
        return recommendationService.listReceivedRecommendations(
                        PageRequest.ofSize(1000),
                        request,
                        username
                )
                .stream()
                .map(this::fromRecommendationDtoToMediaRecommendationDto)
                .toList();
    }

    private Media fromMovieDbToMediaEntity(MovieDb movie) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        Media media = mapper.convertValue(movie, Media.class);
        media.setExternalId(String.valueOf(movie.getId()));
        media.setMediaProvider(com.teasy.CineCircleApi.models.enums.MediaProvider.THE_MOVIE_DATABASE.name());
        media.setPosterUrl(getCompleteImageUrl(movie.getPosterPath()));
        media.setBackdropUrl(getCompleteImageUrl(movie.getBackdropPath()));
        media.setVoteAverage(movie.getVoteAverage());
        media.setVoteCount(movie.getVoteCount());
        media.setRuntime(movie.getRuntime());
        if (media.getCompleted() == null) {
            media.setCompleted(false);
        }
        media.setMediaType(MediaType.MOVIE.name());
        media.setOriginalTitle(movie.getOriginalTitle());
        media.setTitle(movie.getTitle());
        if (movie.getReleaseDate() != null && !movie.getReleaseDate().isEmpty()) {
            media.setReleaseDate(LocalDate.parse(movie.getReleaseDate()));
        }
        media.setOriginalLanguage(movie.getOriginalLanguage());
        if (media.getCompleted() == null) {
            media.setCompleted(false);
        }
        return media;
    }

    private Media fromMTvSeriesToMediaEntity(TvSeries tvSeries) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        Media media = mapper.convertValue(tvSeries, Media.class);
        media.setExternalId(String.valueOf(tvSeries.getId()));
        media.setMediaProvider(com.teasy.CineCircleApi.models.enums.MediaProvider.THE_MOVIE_DATABASE.name());
        media.setPosterUrl(getCompleteImageUrl(tvSeries.getPosterPath()));
        media.setBackdropUrl(getCompleteImageUrl(tvSeries.getBackdropPath()));
        media.setVoteAverage(tvSeries.getVoteAverage());
        media.setVoteCount(tvSeries.getVoteCount());
        media.setRuntime(tvSeries.getEpisodeRuntime() != null && !tvSeries.getEpisodeRuntime().isEmpty() ?
                tvSeries.getEpisodeRuntime().getFirst() : null);
        if (media.getCompleted() == null) {
            media.setCompleted(false);
        }
        media.setMediaType(MediaType.TV_SHOW.name());
        media.setOriginalTitle(tvSeries.getOriginalName());
        media.setTitle(tvSeries.getName());
        if (tvSeries.getFirstAirDate() != null && !tvSeries.getFirstAirDate().isEmpty()) {
            media.setReleaseDate(LocalDate.parse(tvSeries.getFirstAirDate()));
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

    private void completeMedia(Media media) {
        initTmdbApi();
        // get casting
        List<PersonCast> cast;
        List<PersonCrew> crew;
        List<Person> persons;
        List<Genre> genres;
        List<Video> videos;
        if (Objects.equals(media.getMediaType(), MediaType.MOVIE.name())) {
            MovieDb movie = tmdbApi.getMovies()
                    .getMovie(Integer.parseInt(media.getExternalId()),
                            language,
                            TmdbMovies.MovieMethod.credits,
                            TmdbMovies.MovieMethod.videos);
            cast = movie.getCredits().getCast();
            crew = movie.getCredits().getCrew();
            genres = movie.getGenres();
            videos = movie.getVideos();
            persons = new ArrayList<>();
        } else if (Objects.equals(media.getMediaType(), MediaType.TV_SHOW.name())) {
            TvSeries tvSeries = tmdbApi.getTvSeries()
                    .getSeries(Integer.parseInt(media.getExternalId()),
                            language,
                            TmdbTV.TvMethod.credits,
                            TmdbTV.TvMethod.videos);
            cast = tvSeries.getCredits().getCast();
            persons = tvSeries.getCreatedBy();
            genres = tvSeries.getGenres();
            videos = tvSeries.getVideos();
            crew = new ArrayList<>();
        } else {
            return;
        }

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

        // mark media as completed to avoid getting details again later
        media.setCompleted(true);
        mediaRepository.save(media);
    }

    private String getTrailerUrl(List<Video> videos) {
        String youtubeVideoKey = videos
                .stream()
                .filter(video -> Objects.equals(video.getSite(), youtubeVideoSite) &&
                        Objects.equals(video.getType(), trailerVideoType))
                .map(Video::getKey)
                .toList()
                .getFirst();
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

    private <T> T fromMediaEntityToDto(Media media, Class<T> toValueType, String authenticatedUsername) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        var result = mapper.convertValue(media, toValueType);
        addRecommendationRatingFields(result, authenticatedUsername);
        return result;
    }

    private MediaRecommendationDto fromRecommendationDtoToMediaRecommendationDto(RecommendationDto recommendationDto) {
        var mapper = new ObjectMapper()
                .disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
                .registerModule(new JavaTimeModule());
        return mapper.convertValue(recommendationDto, MediaRecommendationDto.class);
    }
}
