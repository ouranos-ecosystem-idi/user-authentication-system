CREATE TABLE public.trades (
    trade_id character varying(256) NOT NULL,
    downstream_operator_id character varying(256) NOT NULL,
    upstream_operator_id character varying(256) NOT NULL,
    downstream_trace_id character varying(256) NOT NULL,
    upstream_trace_id character varying(256),
    trade_date timestamp without time zone,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL
);

COMMENT ON TABLE public.trades IS '取引関係情報テーブル';
COMMENT ON COLUMN public.trades.trade_id IS '取引関係情報識別子';
COMMENT ON COLUMN public.trades.downstream_operator_id IS '川下事業者識別子';
COMMENT ON COLUMN public.trades.upstream_operator_id IS '川上事業者識別子';
COMMENT ON COLUMN public.trades.downstream_trace_id IS '川下部品トレース識別子';
COMMENT ON COLUMN public.trades.upstream_trace_id IS '川上部品トレース識別子';
COMMENT ON COLUMN public.trades.trade_date IS '取引日時';
COMMENT ON COLUMN public.trades.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.trades.created_at IS '作成日時';
COMMENT ON COLUMN public.trades.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.trades.updated_at IS '更新日時';
COMMENT ON COLUMN public.trades.updated_user_id IS '更新ユーザ';



ALTER TABLE ONLY public.trades ADD CONSTRAINT trades_pkey PRIMARY KEY (trade_id);
ALTER TABLE ONLY public.trades ADD CONSTRAINT trades_downstream_operator_id_fkey FOREIGN KEY (downstream_operator_id) REFERENCES public.operators(operator_id);
ALTER TABLE ONLY public.trades ADD CONSTRAINT trades_downstream_trace_id_fkey FOREIGN KEY (downstream_trace_id) REFERENCES public.parts(trace_id);
ALTER TABLE ONLY public.trades ADD CONSTRAINT trades_upstream_operator_id_fkey FOREIGN KEY (upstream_operator_id) REFERENCES public.operators(operator_id);