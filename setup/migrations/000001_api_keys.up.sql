CREATE TABLE public.api_keys (
    id character varying(256) DEFAULT gen_random_uuid() NOT NULL,
    api_key character varying(256) NOT NULL,
    application_name character varying(256) NOT NULL,
    application_attribute character varying(256) NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL
);

COMMENT ON TABLE public.api_keys IS 'APIKEYテーブル';
COMMENT ON COLUMN public.api_keys.id IS 'APIKEYID';
COMMENT ON COLUMN public.api_keys.api_key IS 'APIKEY';
COMMENT ON COLUMN public.api_keys.application_name IS 'アプリケーション名';
COMMENT ON COLUMN public.api_keys.application_attribute IS 'アプリケーション属性(DataSpace,Application,Traceability)';
COMMENT ON COLUMN public.api_keys.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.api_keys.created_at IS '作成日時';
COMMENT ON COLUMN public.api_keys.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.api_keys.updated_at IS '更新日時';
COMMENT ON COLUMN public.api_keys.updated_user_id IS '更新ユーザ';

ALTER TABLE ONLY public.api_keys ADD CONSTRAINT api_key_unique UNIQUE (api_key);
ALTER TABLE ONLY public.api_keys ADD CONSTRAINT api_keys_pkey PRIMARY KEY (id);