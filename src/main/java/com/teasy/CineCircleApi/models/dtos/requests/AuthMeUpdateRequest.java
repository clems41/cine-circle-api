package com.teasy.CineCircleApi.models.dtos.requests;

import jakarta.validation.constraints.NotEmpty;
import lombok.Getter;
import lombok.Setter;

public record AuthMeUpdateRequest (
    @NotEmpty String displayName
){
}
