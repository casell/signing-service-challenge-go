package crypto

import (
	"encoding/pem"
	"testing"
)

const fakeName = "FoobAr"

func TestIsValidAlgorithmTrue(t *testing.T) {
	v := IsValidAlgorithm("rsa")
	if !v {
		t.Fatal("Expected to find the algorithm RSA")
	}
}

func TestIsValidAlgorithmFalse(t *testing.T) {
	v := IsValidAlgorithm(fakeName)
	if v {
		t.Fatal("Expected not to find the algorithm", fakeName)
	}
}

func TestGetValidAlgorithms(t *testing.T) {
	v := GetValidAlgorithms()
	if len(v) == 0 {
		t.Fatal("Expected to find some algorithms", fakeName)
	}
}

func TestFromStringUnknown(t *testing.T) {
	g, err := FromString(fakeName)
	if err == nil {
		t.Fatalf("expecting error, got %v", g)
	}
}
func TestFromStringZero(t *testing.T) {
	g, err := FromString("")
	if err == nil {
		t.Fatalf("expecting error, got %v", g)
	}
}

func TestFromStringRSA(t *testing.T) {
	algo := "rsa"
	fromStringinternal(algo, RSA_ALGORITHM_NAME, t)
}

func TestFromStringECC(t *testing.T) {
	algo := "eCc"
	fromStringinternal(algo, ECC_ALGORITHM_NAME, t)
}

func fromStringinternal(algo, expected string, t *testing.T) {
	g, err := FromString(algo)
	if err != nil {
		t.Fatal(err)
	}
	if g.Algorithm() != expected {
		t.Fatalf("got algorithm %s, expected %s", g.Algorithm(), expected)
	}
}

func TestUnmarshalerBadRSA(t *testing.T) {
	g := &RSAGenerator{}
	_, err := g.Unmarshal([]byte{})
	if err == nil {
		t.Fatalf("expecting error unmarshalling empty bytes")
	}
}

func TestUnmarshalerBadRSAPem(t *testing.T) {
	pbytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PRIVATE_KEY",
		Bytes: []byte("privateKEY"),
	})
	g := &RSAGenerator{}
	_, err := g.Unmarshal(pbytes)
	if err == nil {
		t.Fatalf("expecting error unmarshalling bad key")
	}
}

func TestUnmarshalerBadECC(t *testing.T) {
	g := &ECCGenerator{}
	_, err := g.Unmarshal([]byte{})
	if err == nil {
		t.Fatalf("expecting error unmarshalling empty bytes")
	}
}

func TestUnmarshalerBadECCPem(t *testing.T) {
	pbytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE_KEY",
		Bytes: []byte("privateKeyBytes"),
	})
	g := &ECCGenerator{}
	_, err := g.Unmarshal(pbytes)
	if err == nil {
		t.Fatalf("expecting error unmarshalling bad key")
	}
}

func TestMarshalerRSA(t *testing.T) {
	g := &RSAGenerator{}
	marshalerChecks(t, g)
}

func TestMarshalerECC(t *testing.T) {
	g := &ECCGenerator{}
	marshalerChecks(t, g)
}

func marshalerChecks(t *testing.T, g Generator) {
	algorithmName := g.Algorithm()
	kp, err := g.Generate()
	if err != nil {
		t.Fatalf("error generating %s %v", algorithmName, err)
	}
	_, priv, err := kp.Marshal()
	if err != nil {
		t.Fatalf("error marshalling %s %v", algorithmName, err)
	}

	nkp, err := g.Unmarshal(priv)
	if err != nil {
		t.Fatalf("error unmarshalling %s %v", algorithmName, err)
	}

	if !nkp.Equal(kp) {
		t.Fatalf("error restoring %s %v", algorithmName, err)
	}
}
