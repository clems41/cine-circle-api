package com.teasy.CineCircleApi.models.converters;

import com.teasy.CineCircleApi.models.enums.RecommendationType;
import org.springframework.core.convert.converter.Converter;
import org.springframework.stereotype.Component;

@Component
public class StringToRecommendationTypeConverter implements Converter<String, RecommendationType> {
    @Override
    public RecommendationType convert(String source) {
        return RecommendationType.valueOf(source.toUpperCase());
    }
}
