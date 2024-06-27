-- administrator
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'administrator'),
    (SELECT id FROM permissions WHERE description = 'create_user')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'administrator'),
    (SELECT id FROM permissions WHERE description = 'read_user')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'administrator'),
    (SELECT id FROM permissions WHERE description = 'update_user')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'administrator'),
    (SELECT id FROM permissions WHERE description = 'create_approval')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'administrator'),
    (SELECT id FROM permissions WHERE description = 'read_approval')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'administrator'),
    (SELECT id FROM permissions WHERE description = 'action_approval')
);


-- manager
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'manager'),
    (SELECT id FROM permissions WHERE description = 'read_user')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'manager'),
    (SELECT id FROM permissions WHERE description = 'read_approval')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'manager'),
    (SELECT id FROM permissions WHERE description = 'action_approval')
);

-- staff
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'staff'),
    (SELECT id FROM permissions WHERE description = 'read_user')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'staff'),
    (SELECT id FROM permissions WHERE description = 'create_approval')
);
INSERT INTO role_permission (role_id, permission_id) VALUES (
    (SELECT id FROM roles WHERE name = 'staff'),
    (SELECT id FROM permissions WHERE description = 'read_approval')
);
