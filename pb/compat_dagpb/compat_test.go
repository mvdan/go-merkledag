package merkledag_pb

// mirrored in JavaScript @ https://github.com/ipld/js-dag-pb/blob/master/test/test-compat.js

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ipfs/go-cid"
	dagpb "github.com/ipld/go-codec-dagpb"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

var _ cidlink.Link

var dataZero []byte = make([]byte, 0)
var dataSome []byte = []byte{0, 1, 2, 3, 4}
var cidBytes []byte = []byte{1, 85, 0, 5, 0, 1, 2, 3, 4}
var cidCast, _ = cid.Cast(cidBytes)
var zeroName string = ""
var someName string = "some name"
var zeroTsize int64 = 0
var someTsize int64 = 1010
var largeTsize int64 = 9007199254740991 // JavaScript Number.MAX_SAFE_INTEGER

type testCase struct {
	name          string
	node          ipld.Node
	expectedBytes string
	expectedForm  string
}

func qpMap(fn func(ipld.MapAssembler)) ipld.Node {
	node, err := qp.BuildMap(dagpb.Type.PBNode, -1, fn)
	if err != nil {
		panic(err)
	}
	return node
}

var testCases = []testCase{
	// dagpb does not allow missing Links, though that does not affect the
	// encoding, so we likely don't care.
	// {
	// 	name: "empty",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 	}),
	// },
	//{
	//	name: "Data zero",
	//	node: qpMap(func(am ipld.MapAssembler) {
	//		qp.MapEntry(am, "Data", qp.Bytes(dataZero))
	//	}),
	//	expectedBytes: "0a00",
	//	expectedForm: `{
	//		"Data": ""
	//	}`,
	//},
	//{
	//	name:          "Data some",
	//	node:          dagpb.PBNode{Data: dataSome},
	//	expectedBytes: "0a050001020304",
	//	expectedForm: `{
	//		"Data": "0001020304"
	//	}`,
	//},

	{
		name: "Links zero",
		node: qpMap(func(am ipld.MapAssembler) {
			qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
			}))
		}),
		expectedBytes: "",
		expectedForm:  "{}",
	},
	{
		name: "Data some Links zero",
		node: qpMap(func(am ipld.MapAssembler) {
			qp.MapEntry(am, "Data", qp.Bytes(dataSome))
			qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
			}))
		}),
		expectedBytes: "0a050001020304",
		expectedForm: `{
			"Data": "0001020304"
		}`,
	},
	// TODO: dagpb might need a bit of tweaking to support this.
	// why do we need it in the encoded form?
	//
	// {
	// 	name: "Links empty",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 		qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
	// 			qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
	// 				// qp.MapEntry(am, "Hash", qp.Link(cidlink.Link{}))
	// 				qp.MapEntry(am, "Hash", qp.Link(nil))
	// 			}))
	// 		}))
	// 	}),
	// 	expectedBytes: "1200",
	// 	expectedForm: `{
	// 		"Links": [
	// 		{}
	// 		]
	// 	}`,
	// },
	// {
	// 	name: "Data some Links empty",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 		qp.MapEntry(am, "Data", qp.Bytes(dataSome))
	// 		qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
	// 			qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
	// 				qp.MapEntry(am, "Hash", qp.Link(nil))
	// 			}))
	// 		}))
	// 	}),
	// 	expectedBytes: "12000a050001020304",
	// 	expectedForm: `{
	// 			"Data": "0001020304",
	// 			"Links": [
	// 			{}
	// 			]
	// 		}`,
	// },
	// {
	// 	name: "Links Hash zero",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 		qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
	// 			qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
	// 				qp.MapEntry(am, "Hash", qp.Link(cidlink.Link{}))
	// 			}))
	// 		}))
	// 	}),
	// 	expectedBytes: "12020a00",
	// 	expectedForm: `{
	// 		"Links": [
	// 			{
	// 				"Hash": ""
	// 			}
	// 		]
	// 	}`,
	// },
	// {
	// 	name: "Links Name zero",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 		qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
	// 			qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
	// 				qp.MapEntry(am, "Name", qp.String(""))
	// 			}))
	// 		}))
	// 	}),
	// 	expectedBytes: "12021200",
	// 	expectedForm: `{
	// 		"Links": [
	// 			{
	// 				"Name": ""
	// 			}
	// 		]
	// 	}`,
	// },
	// {
	// 	name: "Links Name some",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 		qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
	// 			qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
	// 				qp.MapEntry(am, "Name", qp.String(someName))
	// 			}))
	// 		}))
	// 	}),
	// 	expectedBytes: "120b1209736f6d65206e616d65",
	// 	expectedForm: `{
	// 		"Links": [
	// 			{
	// 				"Name": "some name"
	// 			}
	// 		]
	// 	}`,
	// },
	// {
	// 	name: "Links Tsize some",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 		qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
	// 			qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
	// 				qp.MapEntry(am, "Tsize", qp.Int(someTsize))
	// 			}))
	// 		}))
	// 	}),
	// 	expectedBytes: "120318f207",
	// 	expectedForm: `{
	// 		"Links": [
	// 			{
	// 				"Tsize": 1010
	// 			}
	// 		]
	// 	}`,
	// },
	// {
	// 	name: "Links Tsize zero",
	// 	node: qpMap(func(am ipld.MapAssembler) {
	// 		qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
	// 			qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
	// 				qp.MapEntry(am, "Tsize", qp.Int(0))
	// 			}))
	// 		}))
	// 	}),
	// 	expectedBytes: "12021800",
	// 	expectedForm: `{
	// 		"Links": [
	// 			{
	// 				"Tsize": 0
	// 			}
	// 		]
	// 	}`,
	// },
	{
		name: "Links Hash some",
		node: qpMap(func(am ipld.MapAssembler) {
			qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
				qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
					qp.MapEntry(am, "Hash", qp.Link(cidlink.Link{Cid: cidCast}))
				}))
			}))
		}),
		expectedBytes: "120b0a09015500050001020304",
		expectedForm: `{
			"Links": [
				{
					"Hash": "015500050001020304"
				}
			]
		}`,
	},
	{
		name: "Links Hash some Name zero",
		node: qpMap(func(am ipld.MapAssembler) {
			qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
				qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
					qp.MapEntry(am, "Hash", qp.Link(cidlink.Link{Cid: cidCast}))
					qp.MapEntry(am, "Name", qp.String(""))
				}))
			}))
		}),
		expectedBytes: "120d0a090155000500010203041200",
		expectedForm: `{
			"Links": [
				{
					"Hash": "015500050001020304",
					"Name": ""
				}
			]
		}`,
	},
	{
		name: "Links Hash some Name some",
		node: qpMap(func(am ipld.MapAssembler) {
			qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
				qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
					qp.MapEntry(am, "Hash", qp.Link(cidlink.Link{Cid: cidCast}))
					qp.MapEntry(am, "Name", qp.String(someName))
				}))
			}))
		}),
		expectedBytes: "12160a090155000500010203041209736f6d65206e616d65",
		expectedForm: `{
			"Links": [
				{
					"Hash": "015500050001020304",
					"Name": "some name"
				}
			]
		}`,
	},
	{
		name: "Links Hash some Tsize zero",
		node: qpMap(func(am ipld.MapAssembler) {
			qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
				qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
					qp.MapEntry(am, "Hash", qp.Link(cidlink.Link{Cid: cidCast}))
					qp.MapEntry(am, "Tsize", qp.Int(0))
				}))
			}))
		}),
		expectedBytes: "120d0a090155000500010203041800",
		expectedForm: `{
			"Links": [
				{
					"Hash": "015500050001020304",
					"Tsize": 0
				}
			]
		}`,
	},
	{
		name: "Links Hash some Tsize some",
		node: qpMap(func(am ipld.MapAssembler) {
			qp.MapEntry(am, "Links", qp.List(-1, func(am ipld.ListAssembler) {
				qp.ListEntry(am, qp.Map(-1, func(am ipld.MapAssembler) {
					qp.MapEntry(am, "Hash", qp.Link(cidlink.Link{Cid: cidCast}))
					qp.MapEntry(am, "Tsize", qp.Int(largeTsize))
				}))
			}))
		}),
		expectedBytes: "12140a0901550005000102030418ffffffffffffff0f",
		expectedForm: `{
			"Links": [
				{
					"Hash": "015500050001020304",
					"Tsize": 9007199254740991
				}
			]
		}`,
	},
}

func TestCompat(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyRoundTrip(t, tc)
		})
	}
}

func verifyRoundTrip(t *testing.T, tc testCase) {
	actualBytes, err := nodeRoundTripToString(t, tc.node)
	if err != nil {
		t.Fatal(err)
	}

	if actualBytes != tc.expectedBytes {
		t.Errorf("Did not match: want=%s got=%s", tc.expectedBytes, actualBytes)
	}
}

func nodeRoundTripToString(t *testing.T, n ipld.Node) (string, error) {
	var buf bytes.Buffer
	err := dagpb.Marshal(n, &buf)
	if err != nil {
		t.Logf("%s", err)
		return "", err
	}
	bs := buf.Bytes()
	t.Logf("[%v]\n", hex.EncodeToString(bs))
	b := dagpb.Type.PBNode.NewBuilder()
	if err := dagpb.Unmarshal(b, bytes.NewReader(bs)); err != nil {
		t.Errorf("%s", err)
	}
	return hex.EncodeToString(bs), nil
}
