package types

import (
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

type JsonNullBool struct {
	sql.NullBool
}

func (jnb *JsonNullBool) UnmarshalJSON(d []byte) error {
	var b *bool
	if err := json.Unmarshal(d, &b); err != nil {
		return err
	}
	if b == nil {
		jnb.Valid = false
		return nil
	}

	jnb.Valid = true
	jnb.Bool = *b
	return nil
}

func (jnb JsonNullBool) MarshalJSON() ([]byte, error) {
	if jnb.Valid {
		return json.Marshal(jnb.Bool)
	}
	return json.Marshal(nil)
}
