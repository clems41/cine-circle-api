package com.teasy.CineCircleApi;

import com.teasy.CineCircleApi.services.externals.mediaProviders.MediaProvider;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.context.ActiveProfiles;

@SpringBootTest
@ActiveProfiles("test")
class CineCircleApiApplicationTests {
	@MockBean
	MediaProvider mediaProvider;

	@Test
	void contextLoads() {
	}

}
