package objects

import (
	"encoding/json"
	"reflect"
)

type Object interface {
	// GetKind returns the type of the objects.
	GetKind() string

	// GetID returns a unique UUID for the objects.
	GetID() string

	// GetName returns the name of the objects. Names are not unique.
	GetName() string

	// SetID sets the ID of the objects.
	SetID(string)

	// SetName sets the name of the objects.
	SetName(string)
}

type Person struct {
	Name string `json:"name"`
	ID string	`json:"id"`
}

func (p *Person) GetKind() string {
	return reflect.TypeOf(p).String()
}

func (p *Person) GetID() string {
	return p.ID
}

func (p *Person) GetName() string {
	return p.Name
}

func (p *Person) SetID(s string) {
	p.ID = s
}

func (p *Person) SetName(s string) {
	p.Name = s
}

func (p Person) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p Person) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &p)
}
