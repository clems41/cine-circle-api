package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Library;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface LibraryRepository extends JpaRepository<Library, UUID> {
    Boolean existsByUser_IdAndMedia_Id(UUID userId, UUID mediaId);
    Optional<Library> findByUser_IdAndMedia_Id(UUID userId, UUID mediaId);
}
