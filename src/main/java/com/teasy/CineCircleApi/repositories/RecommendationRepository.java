package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.Recommendation;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface RecommendationRepository extends JpaRepository<Recommendation, Long> {
    Page<Recommendation> findAllByReceivers_Id(Pageable pageable, Long userId);
    Page<Recommendation> findAllByReceivers_IdAndMedia_Id(Pageable pageable, Long userId, Long mediaId);
}
