package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func asciiletters() string {
	b := strings.Builder{}

	for i := 'A'; i < 'A'+26; i++ {
		b.WriteRune(i)
	}

	for i := 'a'; i < 'a'+26; i++ {
		b.WriteRune(i)
	}

	return b.String()
}

func digits() string {
	b := strings.Builder{}

	for i := range 10 {
		b.WriteString(fmt.Sprintf("%d", i))
	}

	return b.String()
}

func randStr(n int) string {
	b := strings.Builder{}
	p := asciiletters() + digits()

	for range n {
		b.WriteByte(p[rand.Intn(len(p))])
	}

	return b.String()
}

func GenID() string {
	return randStr(30)
}

func fst[T any](ts []T) T {
	if len(ts) == 0 {
		return *new(T)
	}

	return ts[0]
}

func rst[T any](ts []T) []T {
	if len(ts) == 0 {
		return []T{}
	}

	return ts[1:]
}

func fstnrst[T any](ts []T) (T, []T) {
	return fst(ts), rst(ts)
}
