-- 桁数の変更
ALTER TABLE request_status DROP COLUMN response_due_date;
ALTER TABLE request_status DROP COLUMN completed_count;
ALTER TABLE request_status DROP COLUMN completed_count_modified_at;
ALTER TABLE request_status DROP COLUMN trades_count;
ALTER TABLE request_status DROP COLUMN trades_count_modified_at;