package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface RecommendationRepository extends JpaRepository<Recommendation, UUID> {
    Page<Recommendation> findAllByReceivers_Id(Pageable pageable, UUID userId);
    Page<Recommendation> findAllByReceivers_IdAndMedia_Id(Pageable pageable, UUID userId, UUID mediaId);
}
