package persistence

import (
	"github.com/casell/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type operation interface {
	out() chan *result
}

type baseOperation struct {
	outch chan *result
}

func (o *baseOperation) out() chan *result {
	return o.outch
}

type getOperation struct {
	operation
	id uuid.UUID
}

type addOperation struct {
	operation
	data domain.SigningDevice
}

type putOperation struct {
	operation
	data domain.SigningDevice
}

type listOperation struct {
	operation
}

type result struct {
	data domain.SigningDevice
	err  error
}

type MemoryStore struct {
	items   map[uuid.UUID]domain.SigningDevice
	request chan operation
}

func NewMemoryStore() *MemoryStore {
	m := MemoryStore{
		items:   make(map[uuid.UUID]domain.SigningDevice),
		request: make(chan operation),
	}
	go m.start()
	return &m
}

func (m *MemoryStore) start() {
	for r := range m.request {
		switch x := r.(type) {
		case *getOperation:
			x.out() <- &result{
				data: m.items[x.id],
			}
			close(x.out())
		case *addOperation:
			d := x.data
			id := d.ID()
			m.items[id] = d
			x.out() <- &result{
				data: m.items[id],
			}
			close(x.out())
		case *putOperation:
			d := x.data
			id := d.ID()
			_, ok := m.items[id]
			if ok {
				m.items[id] = d
			}
			x.out() <- &result{
				data: m.items[id],
			}
			close(x.out())
		case *listOperation:
			for _, v := range m.items {
				x.out() <- &result{data: v}
			}
			close(x.out())
		}
	}
}

func (m *MemoryStore) Add(x domain.SigningDevice) error {
	out := make(chan *result)
	m.request <- &addOperation{
		operation: &baseOperation{
			outch: out,
		},
		data: x,
	}
	res := <-out

	return res.err
}

func (m *MemoryStore) Put(x domain.SigningDevice) error {
	out := make(chan *result)
	m.request <- &putOperation{
		operation: &baseOperation{
			outch: out,
		},
		data: x,
	}
	res := <-out

	return res.err
}

func (m *MemoryStore) Get(id uuid.UUID) (domain.SigningDevice, error) {
	out := make(chan *result)
	m.request <- &getOperation{
		operation: &baseOperation{
			outch: out,
		},
		id: id,
	}
	res := <-out
	return res.data, res.err
}

func (m *MemoryStore) List() ([]domain.SigningDevice, error) {
	out := make(chan *result)
	m.request <- &listOperation{
		operation: &baseOperation{
			outch: out,
		},
	}

	list := make([]domain.SigningDevice, 0)

	for res := range out {
		if res.err != nil {
			return nil, res.err
		}
		list = append(list, res.data)
	}

	return list, nil
}
