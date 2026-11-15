package main

import (
	"errors"
	"fmt"
	"unsafe"
)

var charClasses = []string{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"0123456789",
	"!@#$%^&*()-_=+[]{};:,.<>?/~",
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

func randomIndex(n int) (int, error) {
	limit := ^uint32(0) - (^uint32(0) % uint32(n))
	var b [4]byte
	for {
		if err := randomBytes(b[:]); err != nil {
			return 0, err
		}
		v := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
		if v < limit {
			return int(v % uint32(n)), nil
		}
	}
}

func Generate(length int) (string, error) {
	if length < len(charClasses) {
		return "", fmt.Errorf("length must be at least %d", len(charClasses))
	}

	var pool string
	for _, c := range charClasses {
		pool += c
	}

	out := make([]byte, 0, length)
	pick := func(set string) error {
		i, err := randomIndex(len(set))
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
		j, err := randomIndex(i + 1)
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
	list := make([]string, 0, count)
	for i := 0; i < count; i++ {
		p, err := Generate(length)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
