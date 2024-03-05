# TraceId

Go package that helps create and pass `TraceId` header among microservices for simple tracing capabilities and log grouping.

The generated `TraceId` value is [UUIDv7](https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#name-uuid-version-7), which lets you infer the time of the trace creation from its value.

## Example

TODO

## Get time from UUIDv7 value

```
$ go run github.com/go-chi/traceid/cmd/traceid 018e0ee7-3605-7d75-b344-01062c6fd8bc
2024-03-05 14:56:57.477 +0100 CET
```

### Create new UUIDv7:
```
$ go run github.com/go-chi/traceid/cmd/traceid
018e0ee7-3605-7d75-b344-01062c6fd8bc
```

## License
[MIT](./LICENSE)