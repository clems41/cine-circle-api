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
    Optional<User> findByUsernameOrEmail(String username, String email);

    Optional<User> findByEmail(String email);

    Optional<User> findByUsername(String username);

    @Query(value = "SELECT u FROM User u " +
            "WHERE (u.username LIKE %?1% " +
            "OR u.displayName LIKE %?1%)" +
            "AND u.id != ?2",
            countQuery = "SELECT COUNT(u) FROM User u " +
                    "WHERE (u.username LIKE %?1% " +
                    "OR u.displayName LIKE %?1%)" +
                    "AND u.id != ?2")
    Page<User> searchUsers(String query, UUID requestUserId, Pageable pageable);

    @Query(value = "select * from users " +
            "where id in (SELECT user_relationships.related_to_user_id FROM user_relationships WHERE user_relationships.user_id = ?1)" +
            "and (?2 is null OR users.username LIKE %?2%)",
            countQuery = "select COUNT(*) from users " +
                    "where id in (SELECT user_relationships.related_to_user_id FROM user_relationships WHERE user_relationships.user_id = ?1)" +
                    "and (?2 is null OR users.username LIKE %?2%)",
            nativeQuery = true)
    Page<User> getRelatedUsers(UUID requestUserId, String query, Pageable pageable);

    @Query(value = "select users.* from users " +
            "LEFT join recommendations ON recommendations.receiver = users.id AND recommendations.sent_by = ?1 " +
            "where users.id in (SELECT user_relationships.related_to_user_id FROM user_relationships WHERE user_relationships.user_id = ?1) " +
            "and (?2 is null OR users.username LIKE %?2%) " +
            "group by users.id " +
            "order by count(recommendations.*) DESC",
            countQuery = "select count(*) from users " +
                    "where users.id in (SELECT user_relationships.related_to_user_id FROM user_relationships WHERE user_relationships.user_id = ?1) " +
                    "and (?2 is null OR users.username LIKE %?2%)",
            nativeQuery = true)
    Page<User> getRelatedUsersWithRecommendationsSentSorting(UUID requestUserId, String query, Pageable pageable);
}