# Simple TODO CRUD API

## run

### run db

```
$ cd db
$ sudo docker-compose up
```

### init db

```
$ go run init.go
```

### run server

```
$ go run main.go
```


## CRUD Reqeusts

### Read

```
curl -X GET http://localhost:8080/todos
```

### Create

```
curl -X POST http://localhost:8080/todos \
-H "Content-Type: application/json" \
-d '{
    "title": "Sample Todo",
    "description": "This is a sample todo item",
    "limited_at": "2024-08-31T23:59:59Z"
}'

```

### Update

```
curl -X PUT http://localhost:8080/todos/1 \
-H "Content-Type: application/json" \
-d '{
    "title": "Updated Todo",
    "description": "This todo item has been updated",
    "limited_at": "2024-09-30T23:59:59Z"
}'
```

### Delete

```
curl -X DELETE http://localhost:8080/todos/1
```
