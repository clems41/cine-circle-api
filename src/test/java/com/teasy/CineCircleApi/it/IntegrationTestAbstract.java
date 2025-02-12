package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.CineCircleApiApplication;
import com.teasy.CineCircleApi.repositories.*;
import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import com.teasy.CineCircleApi.utils.Authenticator;
import com.teasy.CineCircleApi.utils.DummyDataCreator;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.bean.override.mockito.MockitoBean;


@ActiveProfiles("test")
@SpringBootTest(classes = CineCircleApiApplication.class, webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public abstract class IntegrationTestAbstract {
    @MockitoBean
    MediaProvider mediaProvider;

    @LocalServerPort
    protected int port;

    @Autowired
    protected TestRestTemplate restTemplate;

    @Autowired
    protected MediaRepository mediaRepository;

    @Autowired
    protected CircleRepository circleRepository;

    @Autowired
    protected LibraryRepository libraryRepository;

    @Autowired
    protected RecommendationRepository recommendationRepository;

    @Autowired
    protected ErrorRepository errorRepository;

    @Autowired
    protected UserRepository userRepository;
    protected Authenticator authenticator;
    protected DummyDataCreator dummyDataCreator;

    @BeforeEach
    public void setUp() {
        authenticator = new Authenticator(restTemplate, port);
        dummyDataCreator = new DummyDataCreator(userRepository, mediaRepository, recommendationRepository, libraryRepository, circleRepository);
    }
}
