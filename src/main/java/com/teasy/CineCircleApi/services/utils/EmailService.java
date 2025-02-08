package com.teasy.CineCircleApi.services.utils;

import com.teasy.CineCircleApi.models.enums.ErrorMessage;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.utils.SendEmailRequest;
import jakarta.mail.MessagingException;
import jakarta.mail.internet.InternetAddress;
import jakarta.mail.internet.MimeMessage;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.MailException;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.mail.javamail.MimeMailMessage;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;
import java.util.Objects;

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
        try {
            MimeMessage message = emailSender.createMimeMessage();
            MimeMessageHelper helper = new MimeMessageHelper(message, true, "UTF-8");
            helper.setFrom(new InternetAddress(sender, senderName));
            helper.setTo(sendEmailRequest.receiver());
            helper.setSubject(sendEmailRequest.subject());
            helper.setText(getContentFromTemplateAndValues(
                            sendEmailRequest.templateName(),
                            sendEmailRequest.templateValues()),
                    true);
            emailSender.send(message);
        } catch(MailException | MessagingException | UnsupportedEncodingException e) {
            log.error("Error while sending email to {} : e", sendEmailRequest.receiver(), e);
            throw new ExpectedException(ErrorMessage.ERR_EMAILSERVICE_CANNOT_SEND_EMAIL, e);
        }
    }

    private String getContentFromTemplateAndValues(String templateName, Map<String, String> templateValues) throws MessagingException {
        Path templateFilePath = Path.of(String.join("/", mailTemplatesPath, templateName));
        String result = "";
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
