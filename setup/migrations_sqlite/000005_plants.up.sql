CREATE TABLE plants (
    plant_id character varying(256) NOT NULL,
    operator_id character varying(256) NOT NULL,
    plant_name character varying(256) NOT NULL,
    plant_address character varying(256) NOT NULL,
    deleted_at timestamp,
    created_at timestamp NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp NOT NULL,
    updated_user_id text NOT NULL,
    open_plant_id character varying(26) NOT NULL,
    global_plant_id character varying(256),
    PRIMARY KEY (plant_id, operator_id),
    UNIQUE (operator_id, global_plant_id),
    UNIQUE (operator_id, open_plant_id),
    FOREIGN KEY (operator_id) REFERENCES operators(operator_id) ON UPDATE CASCADE ON DELETE CASCADE
);