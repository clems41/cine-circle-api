package com.teasy.CineCircleApi.repositories;
import com.teasy.CineCircleApi.models.entities.Admin;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface AdminRepository extends JpaRepository<Admin, String> {
}
