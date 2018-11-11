# gencf
Generate html form from golang struct and return...


#### Some code is not support


Alias:

```go
type R int
```

```go
type R = int
```

Pointer:

```golang
type A struct{
	a *int
}
```

```golang
type A struct{
	b *B
}
type B struct{
	a *A
}
```
