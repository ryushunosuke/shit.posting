
# Dependencies

- Postgresql 12
- ffmpeg
- go get github.com/gorilla/mux
- go get github.com/lib/pq

# Installation

```console
su postgres
psql -c "CREATE USER shitposting WITH LOGIN PASSWORD 'shitposting' CREATEDB;"
psql -U shitposting -h localhost shitposting
shitposting
create database shitposting;
\c shitposting;
create table items (item jsonb NOT NULL);
```

# Usage

Change the folder part to whichever folder you want to be used in config.json.
If you're on windows and somehow have come to this part, change the `\n` in `ftype.go`'s `Convert()` to `\r\n`

```console
go run .
```
