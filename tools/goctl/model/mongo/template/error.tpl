package model

import (
	"errors"

	"github.com/jialequ/linux-sdk/core/stores/mon"
)

var (
	ErrNotFound        = mon.ErrNotFound
	ErrInvalidObjectId = errors.New("invalid objectId")
)
