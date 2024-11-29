CREATE TABLE operators (
    operator_id character varying(256) NOT NULL,
    operator_name character varying(256) NOT NULL,
    deleted_at timestamp,
    created_at timestamp NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp NOT NULL,
    updated_user_id text NOT NULL,
    operator_address character varying(256) NOT NULL,
    open_operator_id character varying(20) NOT NULL,
    global_operator_id character varying(256),
    PRIMARY KEY (operator_id),
    UNIQUE (global_operator_id),
    UNIQUE (open_operator_id)
);