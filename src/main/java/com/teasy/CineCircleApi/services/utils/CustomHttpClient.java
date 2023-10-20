package com.teasy.CineCircleApi.services.utils;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.teasy.CineCircleApi.models.utils.CustomHttpClientSendRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.http.client.methods.*;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.ssl.SSLContexts;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import javax.net.ssl.SSLContext;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.X509Certificate;
import java.util.StringJoiner;

@Slf4j
@Service
public class CustomHttpClient {

    private CloseableHttpClient httpClient;

    public CustomHttpClient() {
        try {
            httpClient = HttpClientBuilder.create().setSSLContext(buildSslContext()).build();
        } catch (NoSuchAlgorithmException | KeyStoreException | KeyManagementException e) {
            log.error("cannot create custom http client : {}", e.getMessage());
            httpClient = HttpClientBuilder.create().build();
        }
    }

    public <T> T sendRequest(CustomHttpClientSendRequest httpClientSendRequest, Class<T> responseType) {
        // build url with query parameters
        String finalUrl = httpClientSendRequest.getUrl();
        StringJoiner queryParameters = new StringJoiner("&");
        httpClientSendRequest.getQueryParameters().forEach((queryParameterKey, queryParameterValue) ->
                queryParameters.add(String.format("%s=%s", queryParameterKey, URLEncoder.encode(queryParameterValue, StandardCharsets.UTF_8))));
        if (queryParameters.length() > 0) {
            finalUrl = finalUrl.concat("?").concat(queryParameters.toString());
        }

        // build request with specified method
        HttpRequestBase request = prepareRequest(httpClientSendRequest, finalUrl);

        try {
            // send request
            CloseableHttpResponse httpResponse = httpClient.execute(request);

            // read response
            BufferedReader reader = new BufferedReader(new InputStreamReader(
                    httpResponse.getEntity().getContent()));
            StringBuilder sb = new StringBuilder();
            String line;
            while ((line = reader.readLine()) != null) {
                sb.append(line).append("\n");
            }

            // map to specified class
            ObjectMapper mapper = new ObjectMapper()
                    .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
            return mapper.readValue(sb.toString(), responseType);
        } catch (Exception e) {
            throw new ResponseStatusException(
                    HttpStatus.INTERNAL_SERVER_ERROR,
                    String.format("error during sending request %s to %s : %s",
                            httpClientSendRequest.getHttpMethod().toString(),
                            finalUrl,
                            e.getMessage()));
        }
    }

    private static HttpRequestBase prepareRequest(CustomHttpClientSendRequest httpClientSendRequest, String finalUrl) {
        HttpRequestBase request = new HttpGet(finalUrl);
        switch (httpClientSendRequest.getHttpMethod().toString()) {
            case "GET" -> request = new HttpGet(finalUrl);
            case "POST" -> request = new HttpPost(finalUrl);
            case "PUT" -> request = new HttpPut(finalUrl);
            case "DELETE" -> request = new HttpDelete(finalUrl);
        }

        // add header
        if (httpClientSendRequest.getAuthorizationHeader() != null) {
            request.addHeader(HttpHeaders.AUTHORIZATION, httpClientSendRequest.getAuthorizationHeader());
        }
        return request;
    }

    private SSLContext buildSslContext() throws NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        return SSLContexts.custom()
                .loadTrustMaterial(null, (final X509Certificate[] chain, final String authType) -> true)
                .build();
    }
}
