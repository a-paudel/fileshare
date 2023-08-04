-- name: CreateFile :exec
insert into files (code, filename, filesize) values (?, ?, ?);

-- name: GetFileByCode :one
select * from files where code = ?;

-- name: CheckIfCodeExists :one
select count(*) from files where code = ?;

-- name: GetExpiredFiles :many
select * from files where expires_at < current_timestamp;

-- name: DeleteFileByCode :exec
delete from files where code = ?;
