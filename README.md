# decimal.Decimal{} codec

[![Go Doc](https://img.shields.io/badge/go-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/muktihari/decimalcodec)
[![Go Report Card](https://goreportcard.com/badge/github.com/muktihari/decimalcodec)](https://goreportcard.com/report/github.com/muktihari/decimalcodec)

Encoder and Decoder Codec for [decimal.Decimal{}](https://github.com/shopspring/decimal) as [primitive.Decimal128](https://github.com/mongodb/mongo-go-driver/blob/master/bson/primitive/decimal.go) on official [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)

## Usage
```go
import (
    ...
    "github.com/muktihari/decimalcodec"
)

func main() {
    rb := bsoncodec.NewRegistryBuilder()

    // you might want to include defaults encoder and decoder as well
    bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)

    decimalcodec.RegisterEncodeDecoder(rb)
    
    registry := rb.Build()
    
    client, err := mongo.Connect(context.Background(),
        options.Client().
            ApplyURI("mongodb://localhost:27017").
            SetRegistry(registry),
    )
    ...
}

```