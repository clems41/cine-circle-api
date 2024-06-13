package com.teasy.CineCircleApi.repositories;
import com.teasy.CineCircleApi.models.entities.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface UserRepository extends JpaRepository<User, UUID> {
    public Optional<User> findByUsernameOrEmail(String username, String email);
    public Optional<User> findByEmail(String email);
    public Optional<User> findByUsername(String username);

    @Query(value = "SELECT u FROM User u " +
            "WHERE (u.username LIKE %?1% " +
            "OR u.displayName LIKE %?1%)" +
            "AND u.id != ?2",
            countQuery = "SELECT COUNT(u) FROM User u " +
                    "WHERE (u.username LIKE %?1% " +
                    "OR u.displayName LIKE %?1%)" +
                    "AND u.id != ?2")
    Page<User> searchUsers(String query, UUID requestUserId, Pageable pageable);
}
