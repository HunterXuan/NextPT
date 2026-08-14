package mod

import api "server/api/mod"

type ControllerV1 struct{}

func NewV1() api.IModV1 {
	return &ControllerV1{}
}
