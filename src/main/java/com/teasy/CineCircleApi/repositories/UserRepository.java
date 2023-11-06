package com.teasy.CineCircleApi.repositories;
import com.teasy.CineCircleApi.models.entities.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface UserRepository extends JpaRepository<User, UUID> {
    public Optional<User> findByUsernameOrEmail(String username, String email);
    public Optional<User> findByEmail(String email);
    public Optional<User> findByUsername(String username);
}
