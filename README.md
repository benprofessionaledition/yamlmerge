# configgen
Recursively generate configuration files without external dependencies.

Given a file `app.yaml` that looks like this: 
```
base-app-config:
  port: 5000
  endpoint: http://foo.com
  
envs:
  prod:
    config-override:
      endpoint: http://prod.foo.com
      
  test:
    config-override:
      endpoint: http://test.foo.com
```

run `go build main.go --role test --output configurations/` to make a `configurations/test.yaml` that looks like this:

```
port: 5000
endpoint: http://test.foo.com
```

That's it. No Helm, no ktempl, no config servers, no bullshit. 

Other options include `--getroles` to show all roles and exit, and `--input` to specify an alternative input path. 
