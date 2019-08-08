# reko

path based reverse proxy for consul services

## Environment Variables

* `BIND`, bind address, default `127.0.0.1:9001`

## Usage

```
// service with tags
GET /service-name:tag1,tag2/path1/path2

// service with specified id
GET /service-name@service-id/path1/path2
```

## Persisted Storage

`reko` will register self as a `consul` service, the instance id will be persisted to current directory

```
reko-id
```

To deregister current reko instance, run `reko --deregister` manually.

## Credits

Guo Y.K., MIT License
