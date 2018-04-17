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

