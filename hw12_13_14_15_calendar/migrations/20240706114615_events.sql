-- +goose Up
-- +goose StatementBegin
create table if not exists events (
    event_id        bigserial,
    title           varchar(256) not null,
    start_date      timestamp not null,
    end_date        timestamp not null,
    description     text,
    user_id         bigint,
    notify_before   interval,
	constraint events_pk primary key (event_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists events;
-- +goose StatementEnd
