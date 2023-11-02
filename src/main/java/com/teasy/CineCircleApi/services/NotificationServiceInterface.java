package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;

public interface NotificationServiceInterface {
    NotificationTopicResponse getTopicForUsername(String username);
    void sendRecommendation(RecommendationDto recommendation);
}
