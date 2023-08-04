-- +goose Up
-- +goose StatementBegin
create table files (
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp,

    code text primary key,
    filename text not null,
    filesize integer not null,
    expires_at datetime not null generated always as (datetime(created_at, '+1 day')) stored
);

create trigger files_updated_at before update on files for each row begin update files set updated_at = current_timestamp where id = OLD.id; end;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE files;
-- +goose StatementEnd
