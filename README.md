# Bank App

Make new migration script

```migrate create -ext sql -dir migration -seq init```

Migrate database

up : ```migrate -database "postgres://postgres:password@localhost:5432/socmed_ref1?sslmode=disable" -path db/migrations up```

down : ```migrate -database "postgres://postgres:password@localhost:5432/socmed_ref1?sslmode=disable" -path db/migrations down```
