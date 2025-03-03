ALTER TABLE parts ADD COLUMN parts_label_name character varying(256);
ALTER TABLE parts ADD COLUMN parts_add_info1 character varying(256);
ALTER TABLE parts ADD COLUMN parts_add_info2 character varying(256);
ALTER TABLE parts ADD COLUMN parts_add_info3 character varying(256);
COMMENT ON COLUMN parts.parts_label_name IS '部品名称';
COMMENT ON COLUMN parts.parts_add_info1 IS '部品補足情報1';
COMMENT ON COLUMN parts.parts_add_info2 IS '部品補足情報2';
COMMENT ON COLUMN parts.parts_add_info3 IS '部品補足情報3';