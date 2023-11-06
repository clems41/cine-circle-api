package com.teasy.CineCircleApi.models.entities;


import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Getter
@Entity
@Setter
@Table(name = "circles")
@NoArgsConstructor
public class Circle {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Column(name = "is_public", nullable = false)
    private Boolean isPublic;

    @Column(name = "name", nullable = false)
    private String name;

    @Column(name = "description")
    private String description;

    @ManyToMany
    @JoinTable(
            name = "circle_users",
            joinColumns = @JoinColumn(name = "circle_id"),
            inverseJoinColumns = @JoinColumn(name = "user_id"))
    private Set<User> users;

    @ManyToOne(cascade = {CascadeType.PERSIST, CascadeType.MERGE, CascadeType.REFRESH})
    @JoinColumn(name = "created_by", referencedColumnName = "id", nullable=false)
    private User createdBy;

    public Circle(Boolean isPublic,
                  String name,
                  String description,
                  User createdBy) {
        this.isPublic = isPublic;
        this.name = name;
        this.description = description;
        this.createdBy = createdBy;
        this.users = new HashSet<>();
        this.users.add(createdBy);
    }

    public void addUser(User user) {
        this.users.add(user);
    }

    public void removeUser(User user) {
        this.users.remove(user);
    }
}
