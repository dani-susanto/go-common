package main

import (
	"context"
	"errors"

	"github.com/dani-susanto/go-common/log"
)

func main() {
	log := log.New(context.Background(), "common")
	log.Info("oke ini mah ", errors.New("error occured"))
	log.Error("terjadi error", errors.New("error"), "apa lagi nih")
}
