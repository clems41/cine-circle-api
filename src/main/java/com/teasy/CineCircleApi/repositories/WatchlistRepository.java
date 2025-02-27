package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Watchlist;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface WatchlistRepository extends JpaRepository<Watchlist, UUID> {
    Boolean existsByUser_IdAndMedia_Id(UUID userId, UUID mediaId);
    Page<Watchlist> findByUser_Id(UUID userId, Pageable pageable);
    Optional<Watchlist> findByUser_IdAndMedia_Id(UUID userId, UUID mediaId);
}
