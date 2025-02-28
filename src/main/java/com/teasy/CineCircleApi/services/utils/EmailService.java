package com.teasy.CineCircleApi.services.utils;

import com.teasy.CineCircleApi.models.entities.Email;
import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;
import com.teasy.CineCircleApi.models.utils.SendEmailRequest;
import com.teasy.CineCircleApi.repositories.EmailRepository;
import jakarta.mail.MessagingException;
import jakarta.mail.internet.InternetAddress;
import jakarta.mail.internet.MimeMessage;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.MailException;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;

@Service
@Slf4j
public class EmailService {
    private final EmailRepository emailRepository;
    private final JavaMailSender emailSender;

    private final static String sender = "noreply@hucoapp.io";
    private final static String senderName = "HuCo";

    @Value("${spring.mail.templates.path}")
    private String mailTemplatesPath;

    @Autowired
    public EmailService(EmailRepository emailRepository, JavaMailSender emailSender) {
        this.emailRepository = emailRepository;
        this.emailSender = emailSender;
    }

    public void sendEmail(SendEmailRequest sendEmailRequest) throws ExpectedException {
        // keep track of email in database
        var email = new Email();
        email.setSender(sender);
        email.setReceiver(sendEmailRequest.receiver());
        email.setSubject(sendEmailRequest.subject());
        email.setTemplateName(sendEmailRequest.templateName());
        email.setTemplateValues(StringUtils.join(sendEmailRequest.templateValues()));
        emailRepository.save(email);

        // build and send email using MIME
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
        } catch (MessagingException | UnsupportedEncodingException e) {
            updateEmailInDatabase(email, e);
            throw new ExpectedException(ErrorDetails.ERR_EMAIL_BUILDING_REQUEST.addingArgs(email.getId()), e);
        }
        try {
            emailSender.send(message);
            updateEmailInDatabase(email, null);
        } catch(MailException e) {
            updateEmailInDatabase(email, e);
            throw new ExpectedException(ErrorDetails.ERR_EMAIL_SENDING_REQUEST.addingArgs(email.getId(), sendEmailRequest.receiver()), e);
        }
    }

    private void updateEmailInDatabase(Email email, Exception e) {
        if (e != null) {
            email.setError(e.getMessage());
        } else {
            email.setSent(true);
        }
        emailRepository.save(email);
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
