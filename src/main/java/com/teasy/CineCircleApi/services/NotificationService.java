package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.RecommendationDto;
import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.CustomExceptionHandler;
import com.teasy.CineCircleApi.repositories.UserRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.stereotype.Service;

import java.util.Arrays;
import java.util.UUID;

@Service
@Slf4j
public class NotificationService implements NotificationServiceInterface {
    private final SimpMessagingTemplate messagingTemplate;
    private final UserRepository userRepository;

    private final static String topicPrefix = "/topic";

    @Autowired
    public NotificationService(SimpMessagingTemplate messagingTemplate,
                               UserRepository userRepository) {
        this.messagingTemplate = messagingTemplate;
        this.userRepository = userRepository;
    }

    public NotificationTopicResponse getTopicForUsername(String username) {
        var user = userRepository
                .findByUsername(username)
                .orElseThrow(() -> CustomExceptionHandler.userWithUsernameNotFound(username));
        return new NotificationTopicResponse(getTopicFromUser(user));
    }

    public void sendRecommendation(RecommendationDto recommendation) {
        if (recommendation.getReceivers() != null) {
            recommendation.getReceivers().forEach(receiver -> {
                var user = userRepository.findById(UUID.fromString(receiver.getId()));
                if (user.isEmpty()) {
                    log.warn("user with id {} cannot be found and will not received notification", receiver.getId());
                    return;
                }
                messagingTemplate.convertAndSend(getTopicFromUser(user.get()), recommendation);
            });
        }
        if (recommendation.getCircles() != null) {
            recommendation.getCircles().forEach(circle -> circle.getUsers().forEach(receiver -> {
                var user = userRepository.findById(UUID.fromString(receiver.getId()));
                if (user.isEmpty()) {
                    log.warn("user with id {} cannot be found and will not received notification", receiver.getId());
                    return;
                }
                messagingTemplate.convertAndSend(getTopicFromUser(user.get()), recommendation);
            }));
        }
    }

    private String getTopicFromUser(User user) {
        return String.format("%s/%s", topicPrefix, user.getTopicName());
    }
}
