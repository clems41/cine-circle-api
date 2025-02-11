package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.responses.NotificationTopicResponse;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import com.teasy.CineCircleApi.models.entities.User;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.messaging.MessagingException;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.stereotype.Service;

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

    public void sendRecommendation(Recommendation recommendation) {
        try {
            messagingTemplate.convertAndSend(getTopicFromUser(recommendation.getReceiver()), recommendation);
        } catch (MessagingException e) {
            log.error("Error chile sending recommendation through WebSocket for user {} : ",
                    recommendation.getReceiver().getId(), e);
        }
    }
    private String getTopicFromUser(User user) {
        return String.format("%s/%s", topicPrefix, user.getTopicName());
    }
}
