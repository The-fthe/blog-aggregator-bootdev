## bootdev blog-aggregator-bootdev project
- this repo is practice for bootdev blog-aggreagator lesson


### psql generate
```bash 
goose --dir ./sql/schema/ postgres "postgres://<username>:@localhost:5432/blogator" up

```
### queue generate
```bash
sqlc generate
```
[sqlc doc](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
