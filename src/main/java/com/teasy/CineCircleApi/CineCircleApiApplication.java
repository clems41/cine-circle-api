package com.teasy.CineCircleApi;

import io.swagger.v3.oas.annotations.OpenAPIDefinition;
import io.swagger.v3.oas.annotations.enums.SecuritySchemeIn;
import io.swagger.v3.oas.annotations.enums.SecuritySchemeType;
import io.swagger.v3.oas.annotations.info.Info;
import io.swagger.v3.oas.annotations.security.SecurityScheme;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@OpenAPIDefinition(info =
@Info(title = "Cine Circle", version = "2.0.0", description = "La plateforme pour trouver le film du soir !"))
@SecurityScheme(name = "JWT", scheme = "bearer", type = SecuritySchemeType.HTTP, in = SecuritySchemeIn.HEADER)
@SecurityScheme(name = "basic", scheme = "basic", type = SecuritySchemeType.HTTP, in = SecuritySchemeIn.HEADER)
public class CineCircleApiApplication {

	public static void main(String[] args) {
		SpringApplication.run(CineCircleApiApplication.class, args);
	}

}
