package com.teasy.CineCircleApi.repositories;
import com.teasy.CineCircleApi.models.entities.Error;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface ErrorRepository extends JpaRepository<Error, UUID> {
}
