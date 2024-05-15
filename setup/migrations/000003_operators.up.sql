CREATE TABLE public.operators (
    operator_id character varying(256) NOT NULL,
    operator_name character varying(256) NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL,
    operator_address character varying(256) NOT NULL,
    open_operator_id character varying(20) NOT NULL,
    global_operator_id character varying(256)
);

COMMENT ON TABLE public.operators IS '事業者テーブル';
COMMENT ON COLUMN public.operators.operator_id IS '事業者識別子（LEIコード等一意になる文字列を想定）';
COMMENT ON COLUMN public.operators.operator_name IS '企業名';
COMMENT ON COLUMN public.operators.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.operators.created_at IS '作成日時';
COMMENT ON COLUMN public.operators.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.operators.updated_at IS '更新日時';
COMMENT ON COLUMN public.operators.updated_user_id IS '更新ユーザ';
COMMENT ON COLUMN public.operators.operator_address IS '事業者所在地';
COMMENT ON COLUMN public.operators.open_operator_id IS '公開事業者識別子';
COMMENT ON COLUMN public.operators.global_operator_id IS '事業者識別子（グローバル）';

ALTER TABLE ONLY public.operators ADD CONSTRAINT operators_pkey PRIMARY KEY (operator_id);
ALTER TABLE ONLY public.operators ADD CONSTRAINT unique_global_operator_id UNIQUE (global_operator_id);
ALTER TABLE ONLY public.operators ADD CONSTRAINT unique_open_operator_id UNIQUE (open_operator_id);