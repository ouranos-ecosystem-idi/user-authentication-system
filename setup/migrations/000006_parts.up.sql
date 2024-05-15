CREATE TABLE public.parts (
    trace_id character varying(256) NOT NULL,
    operator_id character varying(256) NOT NULL,
    plant_id character varying(256) NOT NULL,
    parts_name character varying(256),
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL,
    support_parts_name character varying(256),
    terminated_flag boolean DEFAULT false NOT NULL,
    amount_required double precision,
    amount_required_unit character varying(256)
);

COMMENT ON TABLE public.parts IS '部品情報テーブル';
COMMENT ON COLUMN public.parts.trace_id IS 'トレース管理識別子';
COMMENT ON COLUMN public.parts.operator_id IS '事業者識別子（外部Key）';
COMMENT ON COLUMN public.parts.plant_id IS '事業所識別子（外部Key）';
COMMENT ON COLUMN public.parts.parts_name IS '部品名';
COMMENT ON COLUMN public.parts.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.parts.created_at IS '作成日時';
COMMENT ON COLUMN public.parts.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.parts.updated_at IS '更新日時';
COMMENT ON COLUMN public.parts.updated_user_id IS '更新ユーザ';
COMMENT ON COLUMN public.parts.support_parts_name IS '補助項目';
COMMENT ON COLUMN public.parts.terminated_flag IS '終端フラグ';
COMMENT ON COLUMN public.parts.amount_required IS '活動量';
COMMENT ON COLUMN public.parts.amount_required_unit IS '活動量単位';

ALTER TABLE ONLY public.parts ADD CONSTRAINT parts_pkey PRIMARY KEY (trace_id);
ALTER TABLE ONLY public.parts ADD CONSTRAINT parts_operator_id_plant_id_fkey FOREIGN KEY (operator_id, plant_id) REFERENCES public.plants(operator_id, plant_id) ON UPDATE CASCADE ON DELETE CASCADE;