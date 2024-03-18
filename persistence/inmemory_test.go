package persistence

import (
	"testing"

	"github.com/casell/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type dummySigningDevice struct {
	id uuid.UUID
}

func (d *dummySigningDevice) Sign(dataToBeSigned string) (string, string, error) {
	panic("unimplemented")
}

func (d *dummySigningDevice) CounterAndLastSignature() (uint, string) {
	panic("unimplemented")
}

func (d *dummySigningDevice) KeyPair() crypto.KeyPair {
	panic("unimplemented")
}

func (d *dummySigningDevice) Label() *string {
	panic("unimplemented")
}

func (d *dummySigningDevice) SignatureAlgorithm() string {
	panic("unimplemented")
}

func (d *dummySigningDevice) ID() uuid.UUID {
	return d.id
}

func TestGetEmpty(t *testing.T) {
	imc := NewMemoryStore()
	uid := uuid.New()
	data, err := imc.Get(uid)
	if err != nil {
		t.Error("Expected nil err, got", err)
	}
	if data != nil {
		t.Fatal("Expected nil err, got", data)
	}
}

func TestPutNotFound(t *testing.T) {
	imc := NewMemoryStore()
	dev := &dummySigningDevice{id: uuid.New()}
	err := imc.Put(dev)
	if err != nil {
		t.Fatal("Expected nil err, got", err)
	}
	r, ok := imc.items[dev.ID()]
	if ok {
		t.Fatalf("Expected to not find item %s, got %v", dev.ID(), r)
	}
}

func TestPutFound(t *testing.T) {
	imc := NewMemoryStore()
	dev := &dummySigningDevice{id: uuid.New()}
	err := imc.Add(dev)
	if err != nil {
		t.Fatal("Expected nil err, got", err)
	}
	r, ok := imc.items[dev.ID()]
	if !ok {
		t.Fatalf("Expected to find item %s, got %v", dev.ID(), imc.items)
	}
	if r != dev {
		t.Fatalf("Expected %v, got %v", dev, r)
	}

}

func TestAdd(t *testing.T) {
	imc := NewMemoryStore()
	dev := &dummySigningDevice{id: uuid.New()}
	err := imc.Add(dev)
	if err != nil {
		t.Fatal("Expected nil err, got", err)
	}
	r, ok := imc.items[dev.ID()]
	if !ok {
		t.Fatalf("Expected to find item %s, got %v", dev.ID(), imc.items)
	}
	if r != dev {
		t.Fatalf("Expected %v, got %v", dev, r)
	}
}

func TestEmptyList(t *testing.T) {
	imc := NewMemoryStore()
	res, err := imc.List()
	if err != nil {
		t.Fatal("Expected nil err, got", err)
	}
	if res == nil {
		t.Fatal("Expected non-nil res, got", res)
	}
	if len(res) != 0 {
		t.Fatal("Expected empty res, got", res)
	}
}

func TestFilledList(t *testing.T) {
	imc := NewMemoryStore()

	dev := &dummySigningDevice{id: uuid.New()}
	err := imc.Add(dev)
	if err != nil {
		t.Fatal("Expected nil PUT err, got", err)
	}
	res, err := imc.List()
	if err != nil {
		t.Fatal("Expected nil err, got", err)
	}
	if res == nil {
		t.Fatal("Expected non-nil res, got", res)
	}
	if len(res) != 1 {
		t.Fatal("Expected res length to be 1, got", res)
	}
	if res[0].ID() != dev.ID() {
		t.Fatalf("Expected res to contain %v, got %v", dev.ID(), res[0])
	}
}
