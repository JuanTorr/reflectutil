# reflectutil

Utility package for reflection functions.

## Instalation

`go get github.com/JuanTorr/reflectutil`

## Functions List

### Simple Values

* *IToBool(i interface{}) (res bool, err error):* Obtains the bool value from the interface

* *IToUint64(i interface{}) (res uint64, err error):* Obtains the uint64 value from the interface

* *IToInt64(i interface{}) (res int64, err error):* Obtains the int64 value from the interface

* *IToFloat64(i interface{}) (res float64, err error):* Obtains the float64 value from the

### Structs

* *TransStructArr(dest, src interface{}) error:* Transforms the source struct array into the destination struct

* *TransStruct(dest, src interface{}) error:* Transforms the source struct into the destination struct

* *MarshalTransStructArr(dest, src interface{}) ([]byte, error):* Marsharls and transforms the source struct array into the destination struct

* *MarshalTransStruct(dest, src interface{}) ([]byte, error):* Marshals and transforms the source struct into the destination struct

## Examples

### IToBool(i interface{}) (res bool, err error)

```go
func main(){
    reflectutil.IToBool("true") //true
    reflectutil.IToBool(1)      //true
    reflectutil.IToBool(5.56)   //true
    reflectutil.IToBool(0)      //false
    reflectutil.IToBool(false)  //false
}
```

### IToUint64(i interface{}) (res uint64, err error)

```go
func main(){
    reflectutil.IToUint64("5")  //5
    reflectutil.IToUint64(6.90) //6
}
```

### IToInt64(i interface{}) (res int64, err error)

```go
func main(){
    reflectutil.IToUint64("5")  //5
    reflectutil.IToUint64(6.90) //6
    reflectutil.IToUint64('R')  //82
}
```

### IToFloat64(i interface{}) (res float64, err error)

```go
func main(){
    reflectutil.IToFloat64("39")   //32
    reflectutil.IToFloat64(int(0)) //32
}
```

### TransStructArr(dest, src interface{}) error

```go
type Src struct{
    A       float64
    B       int
    Bool    int
}

func main(){
    src := []Src{
        Src{A:14.5, B:5, Bool:99},
        Src{A:14.5, B:5, Bool:99},
    }
   type Dest struct{
        Alfa    string  `trans:"from:A"`
        Beta    float64 `trans:"from:B"`
        Bool    bool
    }
    var dest []Dest
    reflectutil.TransStructArr(&dest, src)
}
```

### TransStruct(dest, src interface{}) error

```go
type Src struct{
    A       float64
    B       int
    Bool    int
}

func main(){
    src := Src{A:14.5, B:5, Bool:99}
    dest := struct{
        Alfa    string  `trans:"from:A"`
        Beta    float64 `trans:"from:B"`
        Bool    bool
    }{}

    reflectutil.TransStruct(&dest, src)
    //dest.Alfa:14.5; dest.Beta:5.0; dest.Bool:true
}
```

### MarshalTransStruct(dest, src interface{}) ([]byte, error)

```go

type Src struct{
    A       float64
    B       int
    Bool    int
}

func main(){
    src := Src{A:14.5, B:5, Bool:99}

    json, _ := reflectutil.MarshalTransStruct(struct{
        Alfa    string  `json:"alfa" trans:"from:A"`
        Beta    float64 `json:"beta" trans:"from:B"`
        Bool    bool    `json:"bool"`
    }{}, src)
    //json: {"alfa":"14.5","beta":5,"bool":true}
}
```
