package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) Write(record *Record) error {
	return p.encoder.Encode(&record)
}

func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) Read() (*Record, error) {
	record := &Record{}
	if err := c.decoder.Decode(&record); err != nil {
		return nil, err
	}
	return record, nil
}

func (c *consumer) ReadAll() (map[string]*Record, error) {
	db := make(map[string]*Record)
	for {
		record := &Record{}
		if err := c.decoder.Decode(&record); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		db[record.Key] = record
	}

	return db, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
