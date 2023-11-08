package com.teasy.CineCircleApi.controllers;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.teasy.CineCircleApi.models.dtos.requests.AuthSignUpRequest;
import com.teasy.CineCircleApi.services.TokenService;
import com.teasy.CineCircleApi.services.UserService;
import org.apache.commons.lang3.RandomStringUtils;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(controllers = AuthController.class)
@AutoConfigureMockMvc(addFilters = false)
public class AuthControllerTest {
    @MockBean
    UserService userService;
    @MockBean
    TokenService tokenService;

    @Autowired
    private MockMvc mockMvc;

    @Test
    public void signUp_CheckBodyRequest() throws Exception {
        /* Data */
        var signUpUrl = "/api/v1/auth/sign-up";
        var username = RandomStringUtils.random(10, true, true);
        var tooShortUsername = RandomStringUtils.random(4, true, true);
        var email = String.format("%s@%s.%s",
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(6, true, false),
                RandomStringUtils.random(3, true, false)
        );
        var badEmail1 = RandomStringUtils.random(15, true, false);
        var badEmail2 = String.format("%s@%s",
                RandomStringUtils.random(10, true, false),
                RandomStringUtils.random(6, true, false)
        );
        var password = RandomStringUtils.random(12, true, true);
        var tooShortPassword = RandomStringUtils.random(3, true, true);
        var displayName = RandomStringUtils.random(12, true, false);
        var tooLongDisplayName = RandomStringUtils.random(30, true, false);

        /* Try to sign up with empty username --> should get 400 */
        var requestJson = getSignUpRequestAsJson("", email, password, displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with too short username --> should get 400 */
        requestJson = getSignUpRequestAsJson(tooShortUsername, email, password, displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with empty email --> should get 400 */
        requestJson = getSignUpRequestAsJson(username, "", password, displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with bad formatted email1 --> should get 400 */
        requestJson = getSignUpRequestAsJson(username, badEmail1, password, displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with bad formatted email2 --> should get 400 */
        requestJson = getSignUpRequestAsJson(username, badEmail2, password, displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with empty password --> should get 400 */
        requestJson = getSignUpRequestAsJson(username, email, "", displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with too short password --> should get 400 */
        requestJson = getSignUpRequestAsJson(username, email, tooShortPassword, displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with empty displayName --> should get 400 */
        requestJson = getSignUpRequestAsJson(username, email, password, "");
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with too long displayName --> should get 400 */
        requestJson = getSignUpRequestAsJson(username, email, password, tooLongDisplayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isBadRequest());

        /* Try to sign up with all correct fields --> should get 200 */
        requestJson = getSignUpRequestAsJson(username, email, password, displayName);
        mockMvc.perform(
                        post(signUpUrl)
                                .contentType(MediaType.APPLICATION_JSON)
                                .content(requestJson)
                )
                .andExpect(status().isOk());
    }
    private String getSignUpRequestAsJson(String username, String email, String password, String displayName) throws JsonProcessingException {
        var signUpRequest = new AuthSignUpRequest(username, email, password, displayName);
        ObjectMapper mapper = new ObjectMapper();
        mapper.configure(SerializationFeature.WRAP_ROOT_VALUE, false);
        ObjectWriter ow = mapper.writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(signUpRequest);
    }
}
