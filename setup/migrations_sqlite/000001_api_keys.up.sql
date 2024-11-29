CREATE TABLE api_keys (
    id character varying(256) NOT NULL,
    api_key character varying(256) NOT NULL,
    application_name character varying(256) NOT NULL,
    application_attribute character varying(256) NOT NULL,
    deleted_at timestamp,
    created_at timestamp NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp NOT NULL,
    updated_user_id text NOT NULL,
    UNIQUE(api_key),
    PRIMARY KEY (id)
);