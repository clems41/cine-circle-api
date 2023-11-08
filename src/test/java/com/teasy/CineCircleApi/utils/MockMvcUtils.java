package com.teasy.CineCircleApi.utils;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.assertj.core.api.Assertions;
import org.springframework.test.web.servlet.MvcResult;

import java.io.UnsupportedEncodingException;

public abstract class MockMvcUtils {
    public static <T> T extractResponse(MvcResult result, Class <T> expectedResponseType) throws UnsupportedEncodingException, JsonProcessingException {
        String contentAsString = result.getResponse().getContentAsString();
        var objectMapper = new ObjectMapper();
        return objectMapper.readValue(contentAsString, expectedResponseType);
    }

    public static <T> void extractResponseAndCompare(MvcResult result, T expectedResponse) throws UnsupportedEncodingException, JsonProcessingException {
        var actualResponse = extractResponse(result, expectedResponse.getClass());
        Assertions.assertThat(actualResponse).isEqualTo(expectedResponse);
    }
}
