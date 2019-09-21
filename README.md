# yamlmerge
This is a simple tool that recursively merges two root nodes of a yaml tree and spits the output to stdout. While this isn't exactly rocket science, and can be accomplished 
easily with Python or Perl, this project is written in Go so that it can be compiled to a small standalone binary suitable for use in containers. 

Given a file `app.yaml` that looks like this: 
```
default:
  port: 5000
  endpoint: http://foo.com
  
prod:
  endpoint: http://prod.foo.com
      
test:
  endpoint: http://test.foo.com
```

run `go build main.go --input app.yaml --role test >> configurations/test.yaml` to make a `configurations/test.yaml` that looks like this:

```
port: 5000
endpoint: http://test.foo.com
```

That's it. No Helm, no ktempl, no config servers, no Python, no Perl. 

An exhaustive list of options:
```$xslt

```