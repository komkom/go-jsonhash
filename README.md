## json [ hash ]
Generate persistent MD5 based hashes for maps containig json data. 

Possible usecases
- comparing maps
- check if config files need to be reloaded
- persistable keys for key value stores
- ...

### installing
Use `go get` to install the latest version
of the library.

    > go get -v github.com/komkom/go-jsonhash

to import the package use

```go
import "github.com/komkom/jsonhash"
```

### usage

```go
var j map[string]interface{}
dec := json.NewDecoder(f)
err := dec.Decode(&j); 
if err != nil {
  panic(err)
}

hashValue := jsonhash.Hash(j)
```

### todo
- verify that different json maps have different hashes.
