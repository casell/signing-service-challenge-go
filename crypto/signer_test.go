package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
)

func TestNewGenericSigner(t *testing.T) {
	signer := &rsa.PrivateKey{}
	gs, err := NewGenericSigner(signer, crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	if gs.cryptoHash != crypto.SHA256 {
		t.Fatalf("got hash: %v, expected: %v", gs.cryptoHash, crypto.SHA256)
	}
	if gs.cryptoSigner != signer {
		t.Fatalf("got signer: %v, expected: %v", gs.cryptoSigner, signer)
	}
}

func TestNewGenericSignerNoHASH(t *testing.T) {
	gs, err := NewGenericSigner(&rsa.PrivateKey{}, 99)
	if err == nil {
		t.Fatalf("Expected error got %v", gs)
	}
}

func TestNewGenericSignerSign(t *testing.T) {
	signer, err := (&RSAGenerator{}).Generate()
	if err != nil {
		t.Fatal(err)
	}
	gs, err := NewGenericSigner(signer.PrivateKey(), crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	x := []byte("foobar")
	res, err := gs.Sign(x)

	if err != nil {
		t.Fatal(err)
	}

	hasher := gs.cryptoHash.New()
	_, _ = hasher.Write(x)

	rsakeypair := signer.(*RSAKeyPair)

	if err := rsa.VerifyPKCS1v15(rsakeypair.Public, gs.cryptoHash, hasher.Sum(nil), res); err != nil {
		t.Fatal(err)
	}
}

func TestECC(t *testing.T) {
	pems, _ := base64.StdEncoding.DecodeString("LS0tLS1CRUdJTiBQVUJMSUNfS0VZLS0tLS0KTUhZd0VBWUhLb1pJemowQ0FRWUZLNEVFQUNJRFlnQUU1b1BKMnNHOUp4SndDczNsWEFWWnl3V0JIamt2c291WQp1SHBxS2l6QkVuZGtxV0xFY0hwd1RsU1J4dlBZWXFGbDlyWVN4dTY2V0V0dmxqRDh1Y3ZDbUNjRG8vOWFpa0FnCko2NHFCL0lLdlc2ZUN0TmRTSWFUTGhESXhKaTFPL0FSCi0tLS0tRU5EIFBVQkxJQ19LRVktLS0tLQo=")
	byt, _ := pem.Decode(pems)
	pub, _ := x509.ParsePKIXPublicKey(byt.Bytes)
	signature, _ := base64.StdEncoding.DecodeString("MGUCMQCBo0LSVFAr7WBPc6tt7luxlKq7iaq9qzmjeFbibmuJ9YzAapJuFwQwgHOmmGTOPUcCMGgoldlhBsjM6BVmU4S8JvIkWxeIotywHrhBpkRc1Gf58NfAzaeqVWkyvX3Nct4RXw==")
	signed_data := "3_foo_MGUCMGCW092Sc48PREsC7fnqgUhSyXh5I66h8GwvZaGmxJrWk7u/6d46NmlWZyPtprmUbQIxAPmHke+YbIkvhuV01z4rDc3v0c9AsSxCR4CGVP0+F8jQnZmSlAop5xduGZeMESVu5w=="
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(signed_data))
	hashed := hasher.Sum(nil)
	ret := ecdsa.VerifyASN1(pub.(*ecdsa.PublicKey), hashed, signature)
	if !ret {
		t.Fatal("not valid")
	}
}
