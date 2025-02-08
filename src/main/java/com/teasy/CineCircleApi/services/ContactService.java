package com.teasy.CineCircleApi.services;

import com.teasy.CineCircleApi.models.dtos.requests.ContactSendFeedbackRequest;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.utils.SendEmailRequest;
import com.teasy.CineCircleApi.repositories.AdminRepository;
import com.teasy.CineCircleApi.services.utils.EmailService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class ContactService {
    private final EmailService emailService;
    private final AdminRepository adminRepository;
    private final static String feedbackEmailSubject = "[HuCo] Nouveau retour utilisateur de ";
    private final static String feedbackEmailTemplateName = "contact-feedback.html";
    private final static String feedbackEmailTemplateFeedbackKey = "feedback";

    @Autowired
    public ContactService(EmailService emailService, AdminRepository adminRepository) {
        this.emailService = emailService;
        this.adminRepository = adminRepository;
    }

    public void sendFeedback(ContactSendFeedbackRequest request, String senderUsername) throws ExpectedException {
        var receivers = adminRepository.findAll();
        for (var receiver : receivers) {
            if (receiver.getShouldReceiveFeedback()) {
                Map<String, String> templateValues = new HashMap<>();
                templateValues.put(feedbackEmailTemplateFeedbackKey, request.feedback());
                var sendEmailRequest = new SendEmailRequest(
                        feedbackEmailSubject + senderUsername,
                        receiver.getEmail(),
                        feedbackEmailTemplateName,
                        templateValues
                        );
                emailService.sendEmail(sendEmailRequest);
            }
        }
    }
}
