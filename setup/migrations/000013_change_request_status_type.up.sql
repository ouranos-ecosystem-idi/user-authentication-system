-- 回答系項目追加
ALTER TABLE request_status ADD COLUMN response_due_date timestamp without time zone;
ALTER TABLE request_status ADD COLUMN completed_count numeric;
ALTER TABLE request_status ADD COLUMN completed_count_modified_at timestamp without time zone;
ALTER TABLE request_status ADD COLUMN trades_count numeric;
ALTER TABLE request_status ADD COLUMN trades_count_modified_at timestamp without time zone;
COMMENT ON COLUMN request_status.response_due_date IS '回答希望日';
COMMENT ON COLUMN request_status.completed_count IS '取引関係数';
COMMENT ON COLUMN request_status.completed_count_modified_at IS '取引関係数更新日時';
COMMENT ON COLUMN request_status.trades_count IS '回答完了数';
COMMENT ON COLUMN request_status.trades_count_modified_at IS '回答完了数更新日時';
