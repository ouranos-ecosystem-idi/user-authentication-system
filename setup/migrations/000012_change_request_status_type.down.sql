-- 桁数の変更
ALTER TABLE request_status ALTER COLUMN message TYPE VARCHAR(256);
ALTER TABLE request_status ALTER COLUMN reply_message TYPE VARCHAR(100);