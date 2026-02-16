BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(64) NOT NULL UNIQUE,
    hash VARCHAR(32) NOT NULL,
    auth_salt VARCHAR(16)  NOT NULL,
    encr_salt VARCHAR(16)  NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE users IS 'Stores users for auth.';
COMMENT ON COLUMN users.hash IS 'password hash';
COMMENT ON COLUMN users.auth_salt IS 'salt for password hash';
COMMENT ON COLUMN users.encr_salt IS 'salt for master key';

COMMIT;