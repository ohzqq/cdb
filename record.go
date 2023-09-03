package cdb

import "encoding/json"

type Records struct {
	rows []any
}

func (r Records) Books() ([]Book, error) {
	books := make([]Book, len(r.rows))
	for i, r := range r.rawMsg() {
		err := json.Unmarshal(r, &books[i])
		if err != nil {
			return books, err
		}
	}
	return books, nil
}

func (r Records) StringMap() ([]map[string]any, error) {
	books := make([]map[string]any, len(r.rows))
	for i, r := range r.rawMsg() {
		err := json.Unmarshal(r, &books[i])
		if err != nil {
			return books, err
		}
	}
	return books, nil
}

func (r Records) StringMapString() ([]map[string]string, error) {
	books := make([]map[string]string, len(r.rows))

	anyB, err := r.Books()
	if err != nil {
		return books, err
	}

	for i, b := range anyB {
		books[i] = b.StringMapString()
	}
	return books, nil
}

func (r Records) rawMsg() []json.RawMessage {
	var raw []json.RawMessage
	for _, b := range r.rows {
		raw = append(raw, json.RawMessage(b.(string)))
	}
	return raw
}

func (r Records) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(r.rawMsg())
	if err != nil {
		return []byte{}, err
	}
	return d, nil
}

func (r Records) UnmarshalJSON(d []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	for _, rm := range raw {
		r.rows = append(r.rows, any(rm))
	}

	return nil
}
