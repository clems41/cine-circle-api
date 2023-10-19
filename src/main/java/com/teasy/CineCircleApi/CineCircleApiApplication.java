package com.teasy.CineCircleApi;

import com.teasy.CineCircleApi.config.RsaKeyProperties;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;

@SpringBootApplication
@EnableConfigurationProperties({RsaKeyProperties.class})
public class CineCircleApiApplication {

	public static void main(String[] args) {
		SpringApplication.run(CineCircleApiApplication.class, args);
	}

}
