package com.teasy.CineCircleApi.models.exceptions;

public enum ErrorOnField {
    /* Users fields */
    USER_ID,
    USER_EMAIL,
    USER_USERNAME,
    USER_RESET_PASSWORD_TOKEN,
    USER_PASSWORD,
    USER_DISPLAY_NAME,

    /* Library fields */
    LIBRARY_RATING,

    /* Heading fields */

    /* Contact fields */
    CONTACT_FEEDBACK,

    /* Notification fields */

    /* Recommendation fields */
    RECOMMENDATION_ID,
    RECOMMENDATION_TYPE,
    RECOMMENDATION_RECEIVER,
    RECOMMENDATION_RATING,
    RECOMMENDATION_USER_IDS,
    RECOMMENDATION_MEDIA_ID,

    /* Token fields */

    /* Watchlist fields */

    /* Media fields */
    MEDIA_ID,
    MEDIA_MEDIA_TYPE,

    /* Circle fields */
    CIRCLE_ID,
    CIRCLE_CREATED_BY,
    CIRCLE_NAME,

    /* Email fields */

    /* Internal server fields */

    /* Global fields */
    SEARCH_QUERY,
    UUID,
}
