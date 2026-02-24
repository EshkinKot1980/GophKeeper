BEGIN TRANSACTION;

CREATE TYPE secret_data_type AS ENUM ('credentials', 'card', 'file', 'text');

CREATE TABLE IF NOT EXISTS secrets (
    -- Использются чиловое значение, потому что пользователю удобнее вводить чем UUID
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    data_type secret_data_type NOT NULL,
    name VARCHAR(64) NOT NULL,
    meta_data JSONB NOT NULL DEFAULT '[]'::jsonb,
    encrypted_data BYTEA,
    encrypted_key VARCHAR(128) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_secrets_user_id_data_type ON secrets(user_id, data_type);
COMMIT;
