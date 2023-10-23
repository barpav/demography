package models

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

const MimeTypeNewPersonDataV1 = "application/vnd.newPersonData.v1+json"

// Schema: newPersonData.v1
type NewPersonDataV1 struct {
	Surname    string
	Name       string
	Patronymic string
}

func (m *NewPersonDataV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New person data violates 'newPersonData.v1' schema.")
	}

	m.Surname = strings.TrimSpace(m.Surname)
	m.Name = strings.TrimSpace(m.Name)
	m.Patronymic = strings.TrimSpace(m.Patronymic)

	return m.validate()
}

func (m *NewPersonDataV1) validate() (err error) {
	if m.Surname == "" {
		err = errors.Join(err, errors.New("Person's surname must be specified."))
	}

	if utf8.RuneCountInString(m.Surname) > 150 {
		err = errors.Join(err, errors.New("Person's surname cannot be greater than 150 characters."))
	}

	if m.Name == "" {
		err = errors.Join(err, errors.New("Person's name must be specified."))
	}

	if utf8.RuneCountInString(m.Name) > 150 {
		err = errors.Join(err, errors.New("Person's name cannot be greater than 150 characters."))
	}

	if utf8.RuneCountInString(m.Patronymic) > 150 {
		err = errors.Join(err, errors.New("Person's patronymic cannot be greater than 150 characters."))
	}

	return err
}
