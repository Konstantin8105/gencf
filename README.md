# gencf
Generate html form from golang struct and return...


### Names in HTML form

Name  | Description | Example
--- | --- | ---
Name of struct | First element of name. Point at the end. | `NameOfStruct.` 
Field | Field with Go type, for example : `int`, uint`, `float32`,...  | `Field`
Field | Field with user type(struct). Point at the end | `Field.`
Field [] | Field is slice of Go type, for example : `int`, uint`, `float32`,... In square index of slice | `Field[1]`
Field [] | Field is slice of user type(struct). In square index of slice. Point at the end. | `Field[1].`

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
