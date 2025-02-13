package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.JwtRefreshTokenDto;
import com.teasy.CineCircleApi.models.dtos.JwtTokenDto;
import io.jsonwebtoken.Claims;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.oauth2.jwt.JwtClaimsSet;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.jwt.JwtEncoder;
import org.springframework.security.oauth2.jwt.JwtEncoderParameters;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.Date;
import java.util.Map;
import java.util.UUID;

@Service
public class TokenService {
    private final JwtEncoder encoder;
    private final JwtDecoder decoder;
    private static final Integer hoursBeforeExpirationJwtToken = 24;
    private static final Integer daysBeforeExpirationRefreshToken = 365;

    @Autowired
    public TokenService(JwtEncoder encoder,
                        JwtDecoder decoder) {
        this.encoder = encoder;
        this.decoder = decoder;
    }

    public JwtTokenDto generateToken(String username) {
        Instant now = Instant.now();
        var expirationDate = now.plus(hoursBeforeExpirationJwtToken, ChronoUnit.HOURS);
        JwtClaimsSet claims = JwtClaimsSet.builder()
                .issuer("self")
                .issuedAt(now)
                .expiresAt(expirationDate)
                .subject(username)
                .claim("scope", "")
                .build();
        var tokenString = this.encoder.encode(JwtEncoderParameters.from(claims)).getTokenValue();

        return new JwtTokenDto(tokenString, Date.from(expirationDate));
    }

    public JwtRefreshTokenDto generateRefreshToken() {
        LocalDateTime now = LocalDateTime.now();
        var expirationDate = now.plusDays(daysBeforeExpirationRefreshToken);
        return new JwtRefreshTokenDto(UUID.randomUUID().toString(), expirationDate);
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
