JSON parser in Go using recursive descent parsing.

Used grammar:

```
json   -> object | array | value
object -> '{' [pair (',' pair)*] '}'
pair   -> STRING ':' json
array  -> '[' [json (',' json)*] ']'
value  -> STRING | NUMBER | BOOL | NULL
```

usage:
```sh
go run json-parser file|input <file location | JSON input> <comma-separated list of keys>
```

This repo is not meant for serious use, but is just a thing I did for fun. I covered unicode values, escaping, nested objects, etc, but I haven't used the formal JSON specification. Yet, this parser handles every single input I've tried.

In the end, the parser outputs interface{}, which is either a value, an array, or a map[string]interface{}.

For example, the following JSON can be successfully parsed:
```json
{
  "string": "Hello, world!",
  "number": 12345.6789,
  "booleanTrue": true,
  "booleanFalse": false,
  "nullValue": null,
  "object": {
    "nestedString": "Nested hello",
    "nestedNumber": 42,
    "nestedBoolean": false,
    "nestedArray": [1, "two", 3.0, {"deepObject": "deepValue"}]
  },
  "array": [
    "element1",
    2,
    3.14,
    true,
    null,
    {"arrayObjectKey": "arrayObjectValue"},
    ["nestedArray1", 2, {"nestedArrayObjectKey": "nestedArrayObjectValue"}]
  ],
  "specialCharacters": "!@#$%^&*()_+-=[]{}|;:',.<>/?`~",
  "unicodeCharacters": "こんにちは世界",
  "escapedCharacters": "He said, \"Hello, world!\"",
  "emptyObject": {},
  "emptyArray": []
}
```

License: [MIT](LICENSE)