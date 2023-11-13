package com.teasy.CineCircleApi.mocks;

import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.services.NotificationServiceInterface;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class NotificationServiceMock implements NotificationServiceInterface {
    private Map<String, List<RecommendationDto>> recommendationsSent = new HashMap<>();
    @Override
    public NotificationTopicResponse getTopicForUsername(String username) {
        return new NotificationTopicResponse(username);
    }

    @Override
    public void sendRecommendation(RecommendationDto recommendation) {
        recommendation.getReceivers().forEach(receiver -> {
            if (!recommendationsSent.containsKey(receiver.getUsername())) {
                recommendationsSent.put(receiver.getUsername(), new ArrayList<>());
            }
            var recommendationsForUser = recommendationsSent.get(receiver.getUsername());
            recommendationsForUser.add(recommendation);
            recommendationsSent.put(receiver.getUsername(), recommendationsForUser);
        });
        recommendation.getCircles().forEach(circle -> circle.getUsers().forEach(receiver -> {
            if (!recommendationsSent.containsKey(receiver.getUsername())) {
                recommendationsSent.put(receiver.getUsername(), new ArrayList<>());
            }
            var recommendationsForUser = recommendationsSent.get(receiver.getUsername());
            recommendationsForUser.add(recommendation);
            recommendationsSent.put(receiver.getUsername(), recommendationsForUser);
        }));
    }

    public List<RecommendationDto> getRecommendationsSentForUser(String username) {
        return recommendationsSent.get(username);
    }
}
