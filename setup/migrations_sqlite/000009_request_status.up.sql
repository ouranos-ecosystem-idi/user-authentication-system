CREATE TABLE request_status (
    status_id character varying(256) NOT NULL,
    trade_id character varying(256) NOT NULL,
    request_status text,
    message character varying(256),
    request_type character varying(256) NOT NULL,
    deleted_at timestamp,
    created_at timestamp NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp NOT NULL,
    updated_user_id text NOT NULL,
    reply_message character varying(100),
    cfp_response_status character varying(256),
    trade_tree_status character varying(256),
    PRIMARY KEY (status_id),
    FOREIGN KEY (trade_id) REFERENCES trades(trade_id)
);