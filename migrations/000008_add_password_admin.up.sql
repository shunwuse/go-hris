-- sqlite
INSERT INTO passwords (user_id, hash) VALUES (
    (SELECT id FROM users WHERE username = 'admin'),
    '$2a$10$vGD7tC3abULrbNuwNa//SO5AP72x5iQa84Kw9LYkl6YTNAJDVK1U6'
);
