package com.teasy.CineCircleApi.models.utils;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;
import org.springframework.http.HttpMethod;

import java.util.Map;

@Getter
@Setter
@AllArgsConstructor
public class CustomHttpClientSendRequest {
    private HttpMethod httpMethod;
    private String url;
    private Map<String, String> queryParameters;
    private String authorizationHeader;
}
