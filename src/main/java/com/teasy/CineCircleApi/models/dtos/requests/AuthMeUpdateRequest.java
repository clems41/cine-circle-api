package com.teasy.CineCircleApi.models.dtos.requests;

import org.hibernate.validator.constraints.Length;

public record AuthMeUpdateRequest (
        @Length(min = 4, max = 20, message = "ERR_USER_DISPLAY_NAME_INCORRECT") String displayName
){
}
