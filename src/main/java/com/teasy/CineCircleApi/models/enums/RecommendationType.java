package com.teasy.CineCircleApi.models.enums;

import com.teasy.CineCircleApi.models.exceptions.ErrorDetails;
import com.teasy.CineCircleApi.models.exceptions.ExpectedException;

public enum RecommendationType {
    SENT,
    RECEIVED;

    public static RecommendationType getFromString(String source) throws ExpectedException {
        for(RecommendationType v : values())
            if(v.name().equalsIgnoreCase(source)) return v;
        throw new ExpectedException(ErrorDetails.ERR_RECOMMENDATION_TYPE_NOT_SUPPORTED.addingArgs(source));
    }
}
