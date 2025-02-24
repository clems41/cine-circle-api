package com.teasy.CineCircleApi.services;

import com.nimbusds.jwt.JWTParser;
import com.teasy.CineCircleApi.models.dtos.JwtRefreshTokenDto;
import com.teasy.CineCircleApi.models.dtos.JwtTokenDto;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import io.jsonwebtoken.Claims;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.oauth2.jwt.*;
import org.springframework.stereotype.Service;

import java.text.ParseException;
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

    @Value("${auth.jwt.expiration-delay-in-seconds}")
    private Integer jwtExpirationDelayInSeconds;

    @Value("${auth.refresh-token.expiration-delay-in-days}")
    private Integer refreshTokenExpirationDelayInDays;


    @Autowired
    public TokenService(JwtEncoder encoder,
                        JwtDecoder decoder) {
        this.encoder = encoder;
        this.decoder = decoder;
    }

    public JwtTokenDto generateToken(String username) {
        Instant now = Instant.now();
        var expirationDate = now.plus(jwtExpirationDelayInSeconds, ChronoUnit.SECONDS);
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
        var expirationDate = now.plusDays(refreshTokenExpirationDelayInDays);
        return new JwtRefreshTokenDto(UUID.randomUUID().toString(), expirationDate);
    }

    //retrieve username from jwt token
    public String getUsernameFromToken(String token) {
        return getAllClaimsFromToken(token).get(Claims.SUBJECT).toString();
    }

    //retrieve username from jwt token without checking validity
    public String getUsernameFromTokenWithoutCheckingValidity(String token) throws ExpectedException {
        try {
            return JWTParser.parse(token).getJWTClaimsSet().getSubject();
        } catch (ParseException e) {
            throw new ExpectedException(ErrorDetails.ERR_AUTH_JWT_TOKEN_INVALID.addingArgs(token));
        }
    }

    //for retrieving any information from token we will need the secret key
    private Map<String, Object> getAllClaimsFromToken(String token) {
        return decoder.decode(token).getClaims();
    }
}
