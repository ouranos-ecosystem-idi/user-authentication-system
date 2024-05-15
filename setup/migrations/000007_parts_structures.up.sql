CREATE TABLE public.parts_structures (
    trace_id character varying(256) NOT NULL,
    parent_trace_id character varying(256) NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL
);

COMMENT ON TABLE public.parts_structures IS '部品構成情報テーブル';
COMMENT ON COLUMN public.parts_structures.trace_id IS '部品のトレース識別子';
COMMENT ON COLUMN public.parts_structures.parent_trace_id IS '展開元となる部品のトレース識別子';
COMMENT ON COLUMN public.parts_structures.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.parts_structures.created_at IS '作成日時';
COMMENT ON COLUMN public.parts_structures.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.parts_structures.updated_at IS '更新日時';
COMMENT ON COLUMN public.parts_structures.updated_user_id IS '更新ユーザ';

ALTER TABLE ONLY public.parts_structures ADD CONSTRAINT parts_structures_pkey PRIMARY KEY (trace_id, parent_trace_id);
ALTER TABLE ONLY public.parts_structures ADD CONSTRAINT parts_structures_trace_id_fkey FOREIGN KEY (trace_id) REFERENCES public.parts(trace_id) ON UPDATE CASCADE ON DELETE CASCADE;