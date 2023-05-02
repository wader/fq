package main

import (
	"os"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Writes a BSON document to STDOUT that contains at least one element of every non-deprecated field
// type.

func main() {
	oid, err := primitive.ObjectIDFromHex("644b1619251be740e85522d6")
	if err != nil {
		panic(err)
	}

	d128, err := primitive.ParseDecimal128("123.456")
	if err != nil {
		panic(err)
	}

	b := bsoncore.NewDocumentBuilder().
		AppendDouble("dou", 98.765).
		AppendString("str", "my string").
		AppendDocument("doc", bsoncore.NewDocumentBuilder().
			AppendString("nstr", "nested string").
			AppendArray("narr", bsoncore.NewArrayBuilder().
				AppendInt32(-123).
				Build()).
			AppendDouble("ndou", 98.765).
			Build()).
		AppendArray("arr", bsoncore.NewArrayBuilder().
			AppendString("arr string").
			AppendDocument(bsoncore.NewDocumentBuilder().
				AppendInt32("ni32", -123).
				Build()).
			AppendDouble(98.765).
			Build()).
		AppendBinary("bin", bsontype.BinaryGeneric, []byte{0, 1, 2, 3, 4}).
		AppendObjectID("_id", oid).
		AppendBoolean("boo", true).
		AppendDateTime("dat", 1682642622682).
		AppendNull("nul").
		AppendRegex("reg", "my pattern", "ix").
		AppendJavaScript("jav", "var x = 5;").
		AppendInt32("i32", -123).
		AppendTimestamp("tim", 123, 1682642846).
		AppendInt64("i64", -456).
		AppendDecimal128("dec", d128).
		AppendMinKey("min").
		AppendMaxKey("max").
		Build()

	_, err = os.Stdout.Write(b)
	if err != nil {
		panic(err)
	}
}
