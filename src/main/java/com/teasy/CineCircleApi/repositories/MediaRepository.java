package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Media;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface MediaRepository extends JpaRepository<Media, UUID> {
    @Query(value = "SELECT EXISTS(" +
            "SELECT 1 FROM Library " +
            "WHERE CAST(media.id as string) = ?1 " +
            "AND user.id = (SELECT id FROM User WHERE username = ?2)" +
            ")")
    Boolean isInLibraryOfUser(String mediaId, String username);

    @Query(value = "SELECT EXISTS(" +
            "SELECT 1 FROM Watchlist " +
            "WHERE CAST(media.id as string) = ?1 " +
            "AND user.id = (SELECT id FROM User WHERE username = ?2)" +
            ")")
    Boolean isInWatchlistOfUser(String mediaId, String username);
}
