package store

import (
	"io"
)

type Store interface {
	/*
		Size of the database in bytes on disk
	*/
	Size() int64

	/*
		Writes a snapshot of the database to a writer
	*/
	Backup(io.Writer) error
}
