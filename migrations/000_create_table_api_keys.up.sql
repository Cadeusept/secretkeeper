CREATE TABLE api_keys (
    user_id TEXT NOT NULL,
    service_id TEXT NOT NULL,
    api_key BYTEA NOT NULL,
    PRIMARY KEY (user_id, service_id)
);
