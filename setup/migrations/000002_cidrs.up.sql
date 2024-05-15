CREATE TABLE public.cidrs (
    cidr character varying(18) NOT NULL,
    api_key character varying(256) NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL
);

COMMENT ON TABLE public.cidrs IS 'CIDRテーブル';
COMMENT ON COLUMN public.cidrs.cidr IS 'CIDR';
COMMENT ON COLUMN public.cidrs.api_key IS 'APIキー(外部Key)';
COMMENT ON COLUMN public.cidrs.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.cidrs.created_at IS '作成日時';
COMMENT ON COLUMN public.cidrs.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.cidrs.updated_at IS '更新日時';
COMMENT ON COLUMN public.cidrs.updated_user_id IS '更新ユーザ';

ALTER TABLE ONLY public.cidrs ADD CONSTRAINT cidrs_pkey PRIMARY KEY (cidr, api_key);
ALTER TABLE ONLY public.cidrs ADD CONSTRAINT cidrs_api_key_fkey FOREIGN KEY (api_key) REFERENCES public.api_keys(api_key) ON UPDATE CASCADE ON DELETE CASCADE;