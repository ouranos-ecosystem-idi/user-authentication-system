CREATE TABLE public.apikey_operators (
    api_key character varying(256) NOT NULL,
    operator_id character varying(256) NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    created_user_id text NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    updated_user_id text NOT NULL
);

COMMENT ON COLUMN public.apikey_operators.api_key IS 'APIキー(外部Key)';
COMMENT ON COLUMN public.apikey_operators.operator_id IS '事業者識別子（外部Key）';
COMMENT ON COLUMN public.apikey_operators.deleted_at IS '論理削除日時';
COMMENT ON COLUMN public.apikey_operators.created_at IS '作成日時';
COMMENT ON COLUMN public.apikey_operators.created_user_id IS '作成ユーザ';
COMMENT ON COLUMN public.apikey_operators.updated_at IS '更新日時';
COMMENT ON COLUMN public.apikey_operators.updated_user_id IS '更新ユーザ';

ALTER TABLE ONLY public.apikey_operators ADD CONSTRAINT apikey_operators_pkey PRIMARY KEY (api_key, operator_id);
ALTER TABLE ONLY public.apikey_operators ADD CONSTRAINT apikey_operators_api_key_fkey FOREIGN KEY (api_key) REFERENCES public.api_keys(api_key) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.apikey_operators ADD CONSTRAINT apikey_operators_operator_id_fkey FOREIGN KEY (operator_id) REFERENCES public.operators(operator_id) ON UPDATE CASCADE ON DELETE CASCADE;