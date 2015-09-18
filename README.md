# json2

json2 is an alternative JSON parser.

## Usage

```go
v, err := json2.Unmarshal(`{"foo": true}`)
```

`err` can only be of type `*TokenError` or `*ParseError`.
