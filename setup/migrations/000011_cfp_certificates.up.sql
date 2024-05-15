CREATE TABLE public.cfp_certificates (
    cfp_id character varying(256) NOT NULL,
    id integer NOT NULL,
    cfp_certificate text NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL
);

COMMENT ON TABLE public.cfp_certificates IS 'CFP証明書情報テーブル';
COMMENT ON COLUMN public.cfp_certificates.cfp_id IS 'CFP識別子';
COMMENT ON COLUMN public.cfp_certificates.id IS 'CFP証明書情報識別子';
COMMENT ON COLUMN public.cfp_certificates.cfp_certificate IS '';
COMMENT ON COLUMN public.cfp_certificates.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.cfp_certificates.created_at IS '作成日時';
COMMENT ON COLUMN public.cfp_certificates.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.cfp_certificates.updated_at IS '更新日時';
COMMENT ON COLUMN public.cfp_certificates.updated_user_id IS '更新ユーザ';

ALTER TABLE ONLY public.cfp_certificates ADD CONSTRAINT cfp_certificates_pkey PRIMARY KEY (cfp_id, id);