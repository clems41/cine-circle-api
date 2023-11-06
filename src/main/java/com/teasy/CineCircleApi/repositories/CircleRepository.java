package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Set;
import java.util.UUID;

@Repository
public interface CircleRepository extends JpaRepository<Circle, UUID> {
    List<Circle> findAllByUsers_Id(UUID userId);
}
