package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

var charClasses = []string{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"0123456789",
	"!@#$%^&*()-_=+[]{};:,.<>?/~",
}

var pool = strings.Join(charClasses, "")

const rngBufferSize = 4096

type rng struct {
	buf []byte
	pos int
}

func newRNG() *rng {
	return &rng{buf: make([]byte, rngBufferSize), pos: rngBufferSize}
}

func (r *rng) uint32() (uint32, error) {
	if r.pos+4 > len(r.buf) {
		if err := randomBytes(r.buf); err != nil {
			return 0, err
		}
		r.pos = 0
	}
	v := binary.LittleEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	return v, nil
}

func randomBytes(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	status, _, _ := BCryptGenRandom.Call(
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		bcryptUseSystemPreferredRNG,
	)
	if status != 0 {
		return fmt.Errorf("BCryptGenRandom failed (NTSTATUS 0x%08X)", uint32(status))
	}
	return nil
}

func randomIndex(r *rng, n int) (int, error) {
	limit := ^uint32(0) - (^uint32(0) % uint32(n))
	for {
		v, err := r.uint32()
		if err != nil {
			return 0, err
		}
		if v < limit {
			return int(v % uint32(n)), nil
		}
	}
}

func Generate(length int) (string, error) {
	return generate(newRNG(), length)
}

func generate(r *rng, length int) (string, error) {
	if length < len(charClasses) {
		return "", fmt.Errorf("length must be at least %d", len(charClasses))
	}

	out := make([]byte, 0, length)
	pick := func(set string) error {
		i, err := randomIndex(r, len(set))
		if err != nil {
			return err
		}
		out = append(out, set[i])
		return nil
	}

	for _, set := range charClasses {
		if err := pick(set); err != nil {
			return "", err
		}
	}
	for len(out) < length {
		if err := pick(pool); err != nil {
			return "", err
		}
	}

	for i := len(out) - 1; i > 0; i-- {
		j, err := randomIndex(r, i+1)
		if err != nil {
			return "", err
		}
		out[i], out[j] = out[j], out[i]
	}
	return string(out), nil
}

func GenerateMany(length, count int) ([]string, error) {
	if count < 1 {
		return nil, errors.New("count must be at least 1")
	}
	r := newRNG()
	list := make([]string, 0, count)
	for i := 0; i < count; i++ {
		p, err := generate(r, length)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
