# terraform-provider-gotemplate


## build and run tf
```
go build -o terraform-provider-gotemplate; tf init; tf plan && tf apply
```

## mixed json

when having a mix of json, like
```
{
  "m": "yolo",
  22
}
```

one can use the included `template funcs` to assert the type and change how one deals with the values/keys - example see:
https://gist.github.com/alex-leonhardt/8ed3f78545706d89d466434fb6870023

### template functions

to assert a type when dealing with mixed json, you have the following available:
- isInt
- isString
- isSlice
- isArray
- isMap

and you can use them like this

```
{{ if isInt $v }}
do stuff with {{ $v }}
{{ endif }}

{{ if isMap $v }}
do range over {{ $v }} like ...
{{ range $k, $v := $v -}}
  k={{ $k }}, v={{ $v }}
{{- end }}
{{ endif }}
```

