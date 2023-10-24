package models

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

const MimeTypeEditedPersonDataV1 = "application/vnd.editedPersonData.v1+json"

// Schema: editedPersonData.v1
type EditedPersonDataV1 struct {
	Surname    string
	Name       string
	Patronymic string
	Age        int
	Gender     string
	Country    string
}

func (m *EditedPersonDataV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("Edited person data violates 'editedPersonData.v1' schema.")
	}

	m.Surname = strings.TrimSpace(m.Surname)
	m.Name = strings.TrimSpace(m.Name)
	m.Patronymic = strings.TrimSpace(m.Patronymic)

	return m.validate()
}

func (m *EditedPersonDataV1) validate() (err error) {
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

	if m.Gender != "" && m.Gender != "male" && m.Gender != "female" {
		err = errors.Join(err, errors.New("Incorrect gender value (enum)."))
	}

	return err
}
