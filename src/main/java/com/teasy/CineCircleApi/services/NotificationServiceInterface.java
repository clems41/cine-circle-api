package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;

public interface NotificationServiceInterface {
    NotificationTopicResponse getTopicForUsername(String username) throws ExpectedException;
    void sendRecommendation(RecommendationDto recommendation);
}
