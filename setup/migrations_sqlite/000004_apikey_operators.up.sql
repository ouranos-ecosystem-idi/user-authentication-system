CREATE TABLE apikey_operators (
    api_key character varying(256) NOT NULL,
    operator_id character varying(256) NOT NULL,
    deleted_at timestamp,
    created_at timestamp NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp NOT NULL,
    updated_user_id text NOT NULL,
    PRIMARY KEY (api_key, operator_id),
    FOREIGN KEY (api_key) REFERENCES api_keys(api_key) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES operators(operator_id) ON UPDATE CASCADE ON DELETE CASCADE
);