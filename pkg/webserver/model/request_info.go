package model

import (
	"github.com/google/uuid"
)

type RequestInfo struct {
	RequestId  uuid.UUID
	JSONP      bool
	Callback   string
	Delimiter  bool
	RequestURI string
	Host       string
	RemoteAddr string
}
