package com.teasy.CineCircleApi.models;

import com.teasy.CineCircleApi.utils.validators.UuidValidator;
import com.teasy.CineCircleApi.utils.validators.ValidUuid;
import jakarta.validation.ConstraintValidatorContext;
import jakarta.validation.Valid;
import lombok.Getter;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.commons.lang3.RandomUtils;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;

import java.util.UUID;

public class ValidUuidTest {

    @Mock
    private ConstraintValidatorContext constraintValidatorContext;

    @Test
    public void validUuid() {
        var goodUuid = UUID.randomUUID();
        UuidValidator uuidValidator = new UuidValidator();
        uuidValidator.isValid(goodUuid, constraintValidatorContext);
    }
}
