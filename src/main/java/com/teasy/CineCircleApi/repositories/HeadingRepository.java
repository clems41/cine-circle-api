package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Heading;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface HeadingRepository extends JpaRepository<Heading, UUID> {
    List<Heading> findAllByUserId(UUID userId);
}
