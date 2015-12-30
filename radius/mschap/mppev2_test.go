// tests for MPPE from https://www.ietf.org/rfc/rfc3079.txt
package mschap

import (
	"testing"
	"bytes"
	"fmt"
)

func TestGetMasterKey(t *testing.T) {
	hashHash := []byte{
		0x41, 0xC0, 0x0C, 0x58, 0x4B, 0xD2, 0xD9, 0x1C,
		0x40, 0x17, 0xA2, 0xA1, 0x2F, 0xA5, 0x9F, 0x3F,
	}
	ntRes := []byte {
		0x82, 0x30, 0x9E, 0xCD, 0x8D, 0x70, 0x8B, 0x5E,
		0xA0, 0x8F, 0xAA, 0x39, 0x81, 0xCD, 0x83, 0x54,
		0x42, 0x33, 0x11, 0x4A, 0x3D, 0x85, 0xD6, 0xDF,
	}
	res := getMasterKey(hashHash, ntRes)

	expect := []byte{
		0xFD, 0xEC, 0xE3, 0x71, 0x7A, 0x8C, 0x83, 0x8C,
		0xB3, 0x88, 0xE5, 0x27, 0xAE, 0x3C, 0xDD, 0x31,
	}
	if bytes.Compare(res, expect) != 0 {
		t.Fatal(fmt.Printf("getMasterKey bytes wrong. expect=%d found=%d", expect, res))
	}
}

func TestGetAsymmetricStartKey(t *testing.T) {
	masterKey := []byte{
		0xFD, 0xEC, 0xE3, 0x71, 0x7A, 0x8C, 0x83, 0x8C,
		0xB3, 0x88, 0xE5, 0x27, 0xAE, 0x3C, 0xDD, 0x31,
	}
	expect := []byte{
		0x8B, 0x7C, 0xDC, 0x14, 0x9B, 0x99, 0x3A, 0x1B,
	}
	// Diff 128bit means len=16 and not 8
	res := getAsymmetricStartKey(masterKey, 8, true, true)
	//resRecv := getAsymmetricStartKey(masterKey, 8, false, true)

	if bytes.Compare(res, expect) != 0 {
		t.Fatal(fmt.Printf("GetAsymmetricStartKey bytes wrong. expect=%d found=%d", expect, res))
	}
}

func TestMultipleOfSmaller(t *testing.T) {
	val := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	expect := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	res := multipleOf(val, 16)

	if bytes.Compare(res, expect) != 0 {
		t.Fatal(fmt.Printf("TestMultipleOf bytes wrong. expect=%d found=%d", expect, res))		
	}
}
func TestMultipleOfEqual(t *testing.T) {
	val := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	expect := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	res := multipleOf(val, 16)

	if bytes.Compare(res, expect) != 0 {
		t.Fatal(fmt.Printf("TestMultipleOf bytes wrong. expect=%d found=%d", expect, res))		
	}
}
func TestMultipleOfBigger(t *testing.T) {
	val := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11}
	expect := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	res := multipleOf(val, 16)

	if bytes.Compare(res, expect) != 0 {
		t.Fatal(fmt.Printf("TestMultipleOf bytes wrong. expect=%d found=%d", expect, res))		
	}
}

func TestXor(t *testing.T) {
	a := []byte{0x01}
	b := []byte{0x02}
	expect := []byte{0x03}
	c := xor(a, b)

	if bytes.Compare(c, expect) != 0 {
		t.Fatal(fmt.Printf("TestXor bytes wrong. expect=%d found=%d", expect, c))		
	}
}

func TestMmpe2(t *testing.T) {
	secret := "secret"
	pass := "geheim"
	reqAuth := []byte{}
	ntResponse := []byte{}

	send, recv := Mmpev2(secret, pass, reqAuth, ntResponse)
	if len(send) != 34 {
		t.Fatalf("Send length invalid expect 34, got %d", len(send))
	}
	if len(recv) != 34 {
		t.Fatalf("Recv length invalid expect 34, got %d", len(recv))
	}
}