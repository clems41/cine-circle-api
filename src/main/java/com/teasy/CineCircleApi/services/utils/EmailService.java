package com.teasy.CineCircleApi.services.utils;

import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.utils.SendEmailRequest;
import jakarta.mail.MessagingException;
import jakarta.mail.internet.InternetAddress;
import jakarta.mail.internet.MimeMessage;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.MailException;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;

@Component
@Slf4j
public class EmailService {
    private final JavaMailSender emailSender;

    private final static String sender = "noreply@cinecirlce.com";
    private final static String senderName = "HuCo";

    @Value("${spring.mail.templates.path}")
    private String mailTemplatesPath;

    @Autowired
    public EmailService(JavaMailSender emailSender) {
        this.emailSender = emailSender;
    }

    public void sendEmail(SendEmailRequest sendEmailRequest) throws ExpectedException {
        MimeMessage message = emailSender.createMimeMessage();
        try {
            MimeMessageHelper helper = new MimeMessageHelper(message, true, "UTF-8");
            helper.setFrom(new InternetAddress(sender, senderName));
            helper.setTo(sendEmailRequest.receiver());
            helper.setSubject(sendEmailRequest.subject());
            helper.setText(getContentFromTemplateAndValues(
                            sendEmailRequest.templateName(),
                            sendEmailRequest.templateValues()),
                    true);
        } catch(MessagingException | UnsupportedEncodingException e) {
            throw new ExpectedException(ErrorDetails.ERR_EMAIL_BUILDING_REQUEST, e);
        }
        try {
            emailSender.send(message);
        } catch(MailException e) {
            throw new ExpectedException(ErrorDetails.ERR_EMAIL_SENDING_REQUEST.addingArgs(sendEmailRequest.receiver()), e);
        }
    }

    private String getContentFromTemplateAndValues(String templateName, Map<String, String> templateValues) throws MessagingException {
        Path templateFilePath = Path.of(String.join("/", mailTemplatesPath, templateName));
        String result;
        try {
            result = Files.readString(templateFilePath);
        } catch (IOException e) {
            throw new MessagingException(String.format("cannot get template %s : %s", templateFilePath, e.getMessage()));
        }
        result = templateValues.keySet().stream().reduce(result, (acc, key) -> acc.replaceAll(
                String.format("\\{\\{%s\\}\\}", key),
                templateValues.get(key)));
        return result;
    }
}
