package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Circle;
import com.teasy.CineCircleApi.models.entities.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Set;

public interface CircleRepository extends JpaRepository<Circle, Long> {
    List<Circle> findAllByUsers_Id(Long userId);
}
