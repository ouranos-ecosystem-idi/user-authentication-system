CREATE TABLE cidrs (
    cidr character varying(18) NOT NULL,
    api_key character varying(256) NOT NULL,
    deleted_at timestamp,
    created_at timestamp NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp zone NOT NULL,
    updated_user_id text NOT NULL,
    PRIMARY KEY (cidr, api_key),
    FOREIGN KEY (api_key) REFERENCES api_keys(api_key) ON UPDATE CASCADE ON DELETE CASCADE
);