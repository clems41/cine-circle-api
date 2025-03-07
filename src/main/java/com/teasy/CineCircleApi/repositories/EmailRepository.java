package com.teasy.CineCircleApi.repositories;
import com.teasy.CineCircleApi.models.entities.Email;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface EmailRepository extends JpaRepository<Email, UUID> {
}
