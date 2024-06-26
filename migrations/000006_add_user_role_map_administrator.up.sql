-- sqlite
INSERT INTO user_role (user_id, role_id) VALUES (
    (SELECT id FROM users WHERE username = 'admin'),
    (SELECT id FROM roles WHERE name = 'administrator')
);
