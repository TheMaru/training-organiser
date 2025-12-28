# training-organiser

Private app to organize trainings and keep track of other organisational topics

## Database migrations

Database migrations are handled via [Goose](https://github.com/pressly/goose)

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Migration files are located in folder sql/schema

To migrate to the newest state you can use

```sh
goose postgres postgres://<user>:<host>:<port>/<db-name> up

```
