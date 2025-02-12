package com.teasy.CineCircleApi.it;

import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ErrorResponse;
import com.teasy.CineCircleApi.utils.HttpUtils;
import org.assertj.core.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;

import java.time.LocalDateTime;
import java.util.*;

public class EmailTest extends IntegrationTestAbstract {

    @Test
    public void checkEmailIsStoredInDatabaseWithErrorDetails() {
        /* Data */
        var existingUser = dummyDataCreator.generateUser(true);

        /* Send request to reset password via email, but as the SMTP is not configured for the tests, it should throw an error with cause */
        Map<String, Object> queryParams = new HashMap<>();
        queryParams.put("email", existingUser.getEmail());
        ResponseEntity<ErrorResponse> response = this.restTemplate
                .exchange(
                        HttpUtils.getUriWithQueryParameter(port, HttpUtils.userUrl.concat("reset-password"), queryParams),
                        HttpMethod.GET,
                        new HttpEntity<>(null, null),
                        ErrorResponse.class
                );
        // check response
        var emailsInDatabase = emailRepository.findAll()
                .stream().filter(email -> Objects.equals(email.getReceiver(), existingUser.getEmail())).toList();
        Assertions.assertThat(emailsInDatabase).hasSize(1);
        var email = emailsInDatabase.getFirst();
        var expectedErrorDetails = ErrorDetails.ERR_EMAIL_SENDING_REQUEST.addingArgs(email.getId(), existingUser.getEmail());
        Assertions.assertThat(response.getStatusCode()).isEqualTo(expectedErrorDetails.getHttpStatus());
        Assertions.assertThat(response.getBody()).isNotNull();
        Assertions.assertThat(response.getBody().errorMessage())
                .isEqualTo(expectedErrorDetails.getMessage());
        Assertions.assertThat(response.getBody().errorCode())
                .isEqualTo(expectedErrorDetails.getCode());
        Assertions.assertThat(response.getBody().errorOnObject())
                .isEqualTo(expectedErrorDetails.getErrorOnObject().name());
        Assertions.assertThat(response.getBody().errorOnField()).isNull(); // this field is null for this error
        Assertions.assertThat(response.getBody().errorCause()).isNotNull();
        Assertions.assertThat(response.getBody().errorStack()).isNotNull();

        /* Check that email has been stored in database with the error */
        Assertions.assertThat(email.getReceiver()).isEqualTo(existingUser.getEmail());
        Assertions.assertThat(email.getSender()).isNotEmpty();
        Assertions.assertThat(email.getSubject()).isNotEmpty();
        Assertions.assertThat(email.getTemplateName()).isNotEmpty();
        Assertions.assertThat(email.getTemplateValues()).isNotEmpty();
        Assertions.assertThat(email.getSent()).isFalse();
        Assertions.assertThat(email.getError()).isNotEmpty();
        Assertions.assertThat(email.getSentAt()).isBefore(LocalDateTime.now());
    }
}
