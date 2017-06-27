# K:full_moon:ns - simple data storage library for Golang applications
### Idea - store data in little objects (kons), which can be stored and found by path and every path item is a kone
### Features:
- simple
- flexible
- thread safe
- regexp support for finding
### Examples (see tests):
#### Store object in Kone:
```go
kone := kones.NewKone(nil)
kone.SetData("localhost")
fmt.Println(kone.GetData())
```
output:
> localhost
#### Store object by path:
```go
kone.Upsert(nil, 12, "clients", "john", "balance")
kone.Upsert(nil, 21, "clients", "jessy", "balance")
kone.Upsert(nil, 32, "clients", "bob", "balance")
```
#### Find objects by path (regexp supported):
```go
res, _ := kone.Find(nil, "clients", "[j].*", "balance")
for _, k := range res {
	fmt.Println(k.GetData())
}
```
output:
> 21
> 12
#### Every item is a kone:
```go
bob, _ := kone.FindOne(nil, "clients", "bob")
bob.SetData(45)
fmt.Println(bob.GetData())
```
output:
> 45
#### Pull requests are welcome!
#### Thanks
