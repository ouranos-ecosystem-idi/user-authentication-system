-- 桁数の変更
ALTER TABLE request_status ALTER COLUMN message TYPE VARCHAR(1000);
ALTER TABLE request_status ALTER COLUMN reply_message TYPE VARCHAR(1000);