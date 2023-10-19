package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Media;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;

public interface MediaRepository extends JpaRepository<Media, Long> {
}
