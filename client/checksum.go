package client

import (
	"encoding/base32"
	"errors"

	"github.com/oklog/ulid"
)

func (c *Client) Checksum() (id ulid.ULID, checksum []byte, err error) {
	res, err := c.r.Do("CHECKSUM").Result()
	if err != nil {
		return id, checksum, err
	}

	entries := res.([]interface{})

	sid, ok := entries[0].(string)
	if !ok {
		return id, checksum, errors.New("invalid response")
	}

	id, err = ulid.Parse(sid)
	if err != nil {
		return id, checksum, err
	}

	ssum, ok := entries[1].(string)
	if !ok {
		return id, checksum, errors.New("invalid response")
	}

	checksum, err = base32.StdEncoding.DecodeString(ssum)
	if err != nil {
		return id, checksum, err
	}

	return id, checksum, nil
}
