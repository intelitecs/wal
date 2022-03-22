package log

import (
	"fmt"
	"os"
	"path"

	api "github.com/intelitecs/wal/api/v1/log"
	"google.golang.org/protobuf/proto"
)

type segment struct {
	Store                  *Store
	Index                  *index
	BaseOffset, NextOffset uint64
	Config                 Config
}

func NewSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		BaseOffset: baseOffset,
		Config:     c,
	}
	var err error
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)

	if err != nil {
		return nil, err
	}

	if s.Store, err = NewStore(storeFile); err != nil {
		return nil, err
	}
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE, 0644,
	)
	if err != nil {
		return nil, err
	}
	if s.Index, err = NewIndex(indexFile, c); err != nil {
		return nil, err
	}
	if off, _, err := s.Index.Read(-1); err != nil {
		s.NextOffset = baseOffset
	} else {
		s.BaseOffset = baseOffset + uint64(off) + 1
	}

	return s, nil
}

func (s *segment) Append(record *api.Record) (offset uint64, err error) {
	cur := s.NextOffset
	record.Offset = cur
	p, err := proto.Marshal(record)
	if err != nil {
		return 0, err
	}
	_, pos, err := s.Store.Append(p)
	if err != nil {
		return 0, err
	}

	if err = s.Index.Write(
		uint32(s.NextOffset-uint64(s.BaseOffset)),
		pos,
	); err != nil {
		return 0, err
	}
	s.NextOffset++
	return cur, nil
}

func (s *segment) Read(off uint64) (*api.Record, error) {
	_, pos, err := s.Index.Read(int64(off - s.BaseOffset))
	if err != nil {
		return nil, err
	}
	p, err := s.Store.Read(pos)
	if err != nil {
		return nil, err
	}
	record := &api.Record{}
	err = proto.Unmarshal(p, record)
	return record, err
}

func (s *segment) IsMaxed() bool {
	return s.Store.size >= s.Config.Segment.MaxStoreBytes || s.Index.size >= s.Config.Segment.MaxIndexBytes
}

func (s *segment) Remove() error {
	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.Index.Name()); err != nil {
		return err
	}
	if err := os.Remove(s.Store.Name()); err != nil {
		return err
	}
	return nil
}

func (s *segment) Close() error {
	if err := s.Index.Close(); err != nil {
		return err
	}
	if err := s.Store.Close(); err != nil {
		return err
	}
	return nil
}

func nearestMultiple(j, k uint64) uint64 {
	if j >= 0 {
		return (j / k) * k
	}
	return ((j - k + 1) / k) * k
}
