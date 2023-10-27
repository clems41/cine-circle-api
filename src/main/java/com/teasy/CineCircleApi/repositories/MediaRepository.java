package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Media;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface MediaRepository extends JpaRepository<Media, Long> {
}
