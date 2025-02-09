package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
@Slf4j
public class NotificationService {
    private final SimpMessagingTemplate messagingTemplate;
    private final UserService userService;

    private final static String topicPrefix = "/topic";

    @Autowired
    public NotificationService(SimpMessagingTemplate messagingTemplate, UserService userService) {
        this.messagingTemplate = messagingTemplate;
        this.userService = userService;
    }

    public NotificationTopicResponse getTopicForUsername(String username) throws ExpectedException {
        var user = userService.findUserByUsernameOrElseThrow(username);
        return new NotificationTopicResponse(getTopicFromUser(user));
    }

    public void sendRecommendation(Recommendation recommendation) throws ExpectedException {
        if (recommendation.getReceivers() != null) {
            for (var receiver : recommendation.getReceivers()) {
                sendRecommendationToUser(receiver.getId(), recommendation);
            }
        }
        if (recommendation.getCircles() != null) {
            for (var circle : recommendation.getCircles()) {
                for (var receiver : circle.getUsers()) {
                    sendRecommendationToUser(receiver.getId(), recommendation);
                }
            }
        }
    }

    private void sendRecommendationToUser(UUID userId, Recommendation recommendation) throws ExpectedException {
        User user;
        try {
            user = userService.findUserByIdOrElseThrow(userId);
        } catch (ExpectedException e) {
            if (e.getErrorDetails() != null && e.getErrorDetails().isSameErrorThan(ErrorDetails.ERR_USER_NOT_FOUND)) {
                log.warn("User with id {} cannot be found and will not received notification", userId);
                return;
            } else {
                throw e;
            }
        }
        messagingTemplate.convertAndSend(getTopicFromUser(user), recommendation);
    }

    private String getTopicFromUser(User user) {
        return String.format("%s/%s", topicPrefix, user.getTopicName());
    }
}
