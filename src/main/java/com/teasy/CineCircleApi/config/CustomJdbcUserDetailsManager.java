package com.teasy.CineCircleApi.config;

import org.springframework.jdbc.core.RowMapper;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.provisioning.JdbcUserDetailsManager;

import javax.sql.DataSource;
import java.util.List;

public class CustomJdbcUserDetailsManager extends JdbcUserDetailsManager {
    public CustomJdbcUserDetailsManager(DataSource dataSource) {
        this.setDataSource(dataSource);
    }

    @Override
    protected List<UserDetails> loadUsersByUsername(String username) {
        RowMapper<UserDetails> mapper = (rs, rowNum) -> {
            String username1 = rs.getString(1);
            String password = rs.getString(2);
            boolean enabled = rs.getBoolean(3);
            return new User(username1, password, enabled, true, true, true, AuthorityUtils.NO_AUTHORITIES);
        };
        return this.getJdbcTemplate().query(this.getUsersByUsernameQuery(), mapper, new Object[]{username, username});
    }

    @Override
    protected List<GrantedAuthority> loadUserAuthorities(String username) {
        return this.getJdbcTemplate().query(this.getAuthoritiesByUsernameQuery(), new String[]{username, username}, (rs, rowNum) -> {
            String var10000 = this.getRolePrefix();
            String roleName = var10000 + rs.getString(2);
            return new SimpleGrantedAuthority(roleName);
        });
    }
}
