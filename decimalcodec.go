package decimalcodec

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterEncodeDecoder register decimal.Decimal ValueEncoder and ValueDecoder into *bsoncodec.RegistryBuilder.
func RegisterEncodeDecoder(rb *bsoncodec.RegistryBuilder) {
	var td = reflect.TypeOf(decimal.Decimal{})
	rb.RegisterTypeEncoder(td, bsoncodec.ValueEncoderFunc(DecimalValueEncoder))
	rb.RegisterTypeDecoder(td, bsoncodec.ValueDecoderFunc(DecimalValueDecoder))
}

// DecimalValueEncoder implements bsoncodec.ValueEncoderFunc
func DecimalValueEncoder(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, v reflect.Value) (err error) {
	var d128 primitive.Decimal128

	switch t := v.Interface().(type) {
	case decimal.Decimal:
		d128, err = primitive.ParseDecimal128(t.String())
		if err != nil {
			return err
		}
	default:
		td := reflect.TypeOf(decimal.Decimal{})
		return bsoncodec.ValueEncoderError{
			Name:     "DecimalEncodeValue",
			Types:    []reflect.Type{td, reflect.PtrTo(td)},
			Received: v,
		}
	}
	return vw.WriteDecimal128(d128)
}

// DecimalValueDecoder implements bsoncodec.ValueDecoderFunc
func DecimalValueDecoder(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, v reflect.Value) (err error) {
	var d decimal.Decimal

	switch t := vr.Type(); t {
	case bsontype.Decimal128:
		d128, err := vr.ReadDecimal128()
		if err != nil {
			return err
		}
		d, err = decimal.NewFromString(d128.String())
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil {
			return err
		}
		d = decimal.NewFromFloat(f64)
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil {
			return err
		}
		d = decimal.NewFromInt32(i32)
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil {
			return err
		}
		d = decimal.NewFromInt(i64)
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}
		d, err = decimal.NewFromString(str)
	case bsontype.Null:
		d = decimal.New(0, 0)
	default:
		return fmt.Errorf("received invalid BSON type to decode into decimal.Decimal: %s", vr.Type())
	}

	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(d))

	return
}
