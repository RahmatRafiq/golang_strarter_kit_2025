-- Seeder for default users
INSERT INTO
    users (
        reference,
        username,
        email,
        password,
        jwt_token,
        fcm_token,
        pin,
        created_at,
        updated_at
    )
VALUES
    (
        'USR123456',
        'admin',
        'admin@example.com',
        'argon2id$<hashed-password>',
        '',
        '',
        'argon2id$<hashed-pin>',
        '2025-04-23 22:49:00',
        '2025-04-23 22:49:00'
    );