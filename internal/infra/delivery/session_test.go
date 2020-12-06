package delivery_test

import (
	"foodmap/internal/infra/delivery"
	"testing"
)

func TestSetupSession(t *testing.T) {
	delivery.SetupSession(delivery.NewServer(""), []byte("test secret"))
}
