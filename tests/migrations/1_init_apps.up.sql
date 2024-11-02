INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'sso_secret')
ON CONFLICT DO NOTHING