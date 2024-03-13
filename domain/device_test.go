package domain

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
)

var defaultDeviceFactory = NewDefaultDeviceFactory()

func TestNewDevice(t *testing.T) {
	label := "foo"
	d, err := defaultDeviceFactory.New("RSA", &label)

	if err != nil {
		t.Fatal("unexpected error creating device", err)
	}

	if d.SignatureAlgorithm() != "RSA" {
		t.Fatalf("expected device algorithm to be: RSA, got %v", d.SignatureAlgorithm())
	}

	if d.Label() != &label {
		t.Fatalf("expected device label to be: %v, got %v", label, d.Label())
	}

	counter, lastSignature := d.CounterAndLastSignature()

	if counter != 0 {
		t.Fatalf("expected device counter to be 0, got %v", counter)
	}

	uuid := d.ID().String()
	lastSig, err := base64.StdEncoding.DecodeString(lastSignature)

	if err != nil {
		t.Fatal("unexpected error decoding last signature", err)
	}

	if string(lastSig) != uuid {
		t.Fatalf("expected device last signature to be base64(%s), got base64(%s)", uuid, lastSig)
	}
}

func TestNewDeviceInvalidAlgo(t *testing.T) {
	d, err := defaultDeviceFactory.New("foo", nil)
	if err == nil {
		t.Fatalf("expected invalid algorithm error, got %v", d)
	}
}

func TestRestoreDeviceInvalidAlgo(t *testing.T) {
	d, err := defaultDeviceFactory.Restore(uuid.UUID{}, "foo", nil, 0, "", nil, nil)
	if err == nil {
		t.Fatalf("expected invalid algorithm error, got %v", d)
	}
}

func TestRestoreDevice(t *testing.T) {
	label := "foo"
	d, err := defaultDeviceFactory.New("RSA", &label)
	if err != nil {
		t.Fatal("unexpected error creating device", err)
	}
	counter, lastSignatureB64 := d.CounterAndLastSignature()

	rd, err := defaultDeviceFactory.Restore(
		d.ID(),
		d.SignatureAlgorithm(),
		d.Label(),
		counter,
		lastSignatureB64,
		d.KeyPair(),
		&sync.RWMutex{},
	)

	if err != nil {
		t.Fatal("unexpected error restoring device", err)
	}

	if d.SignatureAlgorithm() != rd.SignatureAlgorithm() {
		t.Fatalf("expected device algorithm to be: %v, got %v", d.SignatureAlgorithm(), rd.SignatureAlgorithm())
	}

	if d.Label() != rd.Label() {
		t.Fatalf("expected device label to be: %v, got %v", d.Label(), rd.Label())
	}

	rcounter, rlastSignatureB64 := rd.CounterAndLastSignature()

	if counter != rcounter {
		t.Fatalf("expected device counter to be %v, got %v", counter, rcounter)
	}

	if lastSignatureB64 != rlastSignatureB64 {
		t.Fatalf("expected device last signature to be %s, got %s", lastSignatureB64, rlastSignatureB64)
	}

	if !rd.KeyPair().Equal(d.KeyPair()) {
		t.Fatalf("expected device keypair to be %s, got %s", d.KeyPair(), rd.KeyPair())
	}

}

type signedDataMetadata struct {
	counter       uint64
	data          string
	lastSignature string
	signature     string
	err           error
}

func getSignedDataMetadata(signature, s string) (*signedDataMetadata, error) {
	metaSlice := strings.Split(s, "_")
	counter, err := strconv.ParseUint(metaSlice[0], 0, 0)
	if err != nil {
		return nil, err
	}
	data := strings.Join(metaSlice[1:len(metaSlice)-1], "_")
	return &signedDataMetadata{
		counter:       counter,
		data:          data,
		lastSignature: metaSlice[len(metaSlice)-1],
		signature:     signature,
	}, nil
}

func TestSign(t *testing.T) {
	d, err := defaultDeviceFactory.New("RSA", nil)
	if err != nil {
		t.Fatal("unexpected error creating device", err)
	}

	initialCounter, initialSignature := d.CounterAndLastSignature()

	requestsNo := 10

	out := make(chan *signedDataMetadata, requestsNo)

	wg := &sync.WaitGroup{}
	wg.Add(requestsNo)

	for i := 0; i < requestsNo; i++ {
		go func(i int, wg *sync.WaitGroup, out chan *signedDataMetadata) {
			signature, dataToBeSigned, err := d.Sign("test")
			if err != nil {
				out <- &signedDataMetadata{err: fmt.Errorf("iteration %d: %v", i, err)}
			}
			hashalgo := crypto.SHA256
			devicepublickey := d.KeyPair().PublicKey()
			if err := verifyRSA(hashalgo, dataToBeSigned, devicepublickey, signature); err != nil {
				out <- &signedDataMetadata{err: fmt.Errorf("iteration %d: %v", i, err)}
			}
			meta, err := getSignedDataMetadata(signature, dataToBeSigned)
			if err != nil {
				out <- &signedDataMetadata{err: fmt.Errorf("iteration %d: %v", i, err)}
			}
			out <- meta
			wg.Done()
		}(i, wg, out)
	}

	wg.Wait()
	close(out)

	replies := make([]*signedDataMetadata, 0, requestsNo)

	for m := range out {
		if m.err != nil {
			t.Fatal(err)
		}
		replies = append(replies, m)
	}

	slices.SortStableFunc(replies, func(a *signedDataMetadata, b *signedDataMetadata) int {
		return int(a.counter - b.counter)
	})

	for i, m := range replies {
		if i == 0 {
			if m.counter != uint64(initialCounter) {
				t.Fatalf("Iteration %d: counter is not greater than previous. current: %v, previous: %v", i, m, initialCounter)
			}
			if m.lastSignature != initialSignature {
				t.Fatalf("Iteration %d: signature do not match. current: %v, previous %v", i, m, initialSignature)
			}
			continue
		}
		prevM := replies[i-1]
		if m.counter <= prevM.counter {
			t.Fatalf("Iteration %d: counter is not greater than previous. current: %v, previous: %v", i, m, prevM)
		}
		if m.lastSignature != prevM.signature {
			t.Fatalf("Iteration %d: signature do not match. current: %v, previous %v", i, m, prevM)
		}
	}
}

func verifyRSA(hashalgo crypto.Hash, dataToBeSigned string, devicepublickey crypto.PublicKey, lastSignature string) error {
	hasher := hashalgo.New()
	_, _ = hasher.Write([]byte(dataToBeSigned))

	publicKey := devicepublickey.(*rsa.PublicKey)

	lastSignatureBytes, err := base64.StdEncoding.DecodeString(lastSignature)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, hashalgo, hasher.Sum(nil), lastSignatureBytes)
}
