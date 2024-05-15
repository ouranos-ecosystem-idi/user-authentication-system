CREATE TABLE public.cfp_infomation (
    trace_id character varying(256) NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL,
    cfp_id character varying(256) NOT NULL,
    ghg_emission numeric,
    ghg_declared_unit character varying(20),
    cfp_type character varying(20) NOT NULL,
    dqr_type character varying(256) NOT NULL,
    te_r numeric,
    ge_r numeric,
    ti_r numeric
);

COMMENT ON TABLE public.cfp_infomation IS 'cfp情報テーブル';
COMMENT ON COLUMN public.cfp_infomation.trace_id IS 'トレース識別子';
COMMENT ON COLUMN public.cfp_infomation.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.cfp_infomation.created_at IS '作成日時';
COMMENT ON COLUMN public.cfp_infomation.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.cfp_infomation.updated_at IS '更新日時';
COMMENT ON COLUMN public.cfp_infomation.updated_user_id IS '更新ユーザ';
COMMENT ON COLUMN public.cfp_infomation.cfp_id IS 'CFP識別子';
COMMENT ON COLUMN public.cfp_infomation.ghg_emission IS 'GHG排出量';
COMMENT ON COLUMN public.cfp_infomation.ghg_declared_unit IS 'GHG排出量の単位';
COMMENT ON COLUMN public.cfp_infomation.cfp_type IS 'CFPの種類';
COMMENT ON COLUMN public.cfp_infomation.dqr_type IS 'DQRの種類';
COMMENT ON COLUMN public.cfp_infomation.te_r IS 'TE_R';
COMMENT ON COLUMN public.cfp_infomation.ge_r IS 'GE_R';
COMMENT ON COLUMN public.cfp_infomation.ti_r IS 'TI_R';

ALTER TABLE ONLY public.cfp_infomation ADD CONSTRAINT cfp_infomation_pkey PRIMARY KEY (trace_id, cfp_type);
ALTER TABLE ONLY public.cfp_infomation ADD CONSTRAINT cfp_infomation_trace_id_fkey FOREIGN KEY (trace_id) REFERENCES public.parts(trace_id) ON UPDATE CASCADE ON DELETE CASCADE;