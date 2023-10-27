package com.teasy.CineCircleApi.repositories;
import com.teasy.CineCircleApi.models.entities.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    public Optional<User> findByUsernameOrEmail(String username, String email);
    public Optional<User> findByEmail(String email);
    public Optional<User> findByUsername(String username);
//    public Optional<User> findById(Long id);
}
