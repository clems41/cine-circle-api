package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.JwtTokenDto;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.jwt.*;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.stereotype.Service;
import org.springframework.security.core.GrantedAuthority;

import java.time.*;
import java.time.temporal.ChronoUnit;
import java.util.Date;
import java.util.Map;
import java.util.Objects;
import java.util.function.Function;
import java.util.stream.Collectors;
import io.jsonwebtoken.Claims;
import org.springframework.web.server.ResponseStatusException;

@Service
public class TokenService {
    private final JwtEncoder encoder;
    private JwtDecoder decoder;
    private static final Integer hoursBeforeExpiration = 24;

    @Autowired
    public TokenService(JwtEncoder encoder,
                        JwtDecoder decoder) {
        this.encoder = encoder;
        this.decoder = decoder;
    }

    public JwtTokenDto generateToken(Authentication authentication) {
        Instant now = Instant.now();
        var expirationDate = now.plus(hoursBeforeExpiration, ChronoUnit.HOURS);
        String scope = authentication.getAuthorities().stream()
                .map(GrantedAuthority::getAuthority)
                .collect(Collectors.joining(" "));
        JwtClaimsSet claims = JwtClaimsSet.builder()
                .issuer("self")
                .issuedAt(now)
                .expiresAt(expirationDate)
                .subject(authentication.getName())
                .claim("scope", scope)
                .build();
        var tokenString = this.encoder.encode(JwtEncoderParameters.from(claims)).getTokenValue();

        return new JwtTokenDto(tokenString, Date.from(expirationDate));
    }

    //retrieve username from jwt token
    public String getUsernameFromToken(String token) {
        return getAllClaimsFromToken(token).get(Claims.SUBJECT).toString();
    }

    //for retrieving any information from token we will need the secret key
    private Map<String, Object> getAllClaimsFromToken(String token) {
        return decoder.decode(token).getClaims();
    }
}
