package com.teasy.CineCircleApi.repositories;

import com.teasy.CineCircleApi.models.entities.Library;
import com.teasy.CineCircleApi.models.entities.Watchlist;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface LibraryRepository extends JpaRepository<Library, Long> {
}
