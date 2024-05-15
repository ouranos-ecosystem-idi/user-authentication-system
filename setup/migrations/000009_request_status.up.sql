CREATE TABLE public.request_status (
    status_id character varying(256) NOT NULL,
    trade_id character varying(256) NOT NULL,
    request_status text,
    message character varying(256),
    request_type character varying(256) NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL,
    reply_message character varying(100),
    cfp_response_status character varying(256),
    trade_tree_status character varying(256)
);

COMMENT ON TABLE public.request_status IS '依頼情報テーブル';
COMMENT ON COLUMN public.request_status.status_id IS '依頼情報識別子';
COMMENT ON COLUMN public.request_status.trade_id IS '取引関係情報識別子';
COMMENT ON COLUMN public.request_status.request_status IS 'リクエストステータス';
COMMENT ON COLUMN public.request_status.message IS 'メッセージ';
COMMENT ON COLUMN public.request_status.request_type IS 'リクエストタイプ';
COMMENT ON COLUMN public.request_status.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.request_status.created_at IS '作成日時';
COMMENT ON COLUMN public.request_status.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.request_status.updated_at IS '更新日時';
COMMENT ON COLUMN public.request_status.updated_user_id IS '更新ユーザ';
COMMENT ON COLUMN public.request_status.reply_message IS '差戻メッセージ';
COMMENT ON COLUMN public.request_status.cfp_response_status IS 'CFPの回答状況';
COMMENT ON COLUMN public.request_status.trade_tree_status IS '取引関係情報終端状況';

ALTER TABLE ONLY public.request_status ADD CONSTRAINT request_status_pkey PRIMARY KEY (status_id);
ALTER TABLE ONLY public.request_status ADD CONSTRAINT request_status_trade_id_fkey FOREIGN KEY (trade_id) REFERENCES public.trades(trade_id);