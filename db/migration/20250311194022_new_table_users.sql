-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS public.users (
    id GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN public.users.id IS 'Идентификатор пользователя';
COMMENT ON COLUMN public.users."name" IS 'Наименование пользователя';
COMMENT ON COLUMN public.users.email IS 'Электронный адрес пользователя';
COMMENT ON COLUMN public.users.created_at IS 'Пользователь создан';
COMMENT ON COLUMN public.users.updated_at IS 'Пользователь изменён';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
