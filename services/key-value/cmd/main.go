package main

import (
	"cloud-native/services/key-value/internal/app/adapters/common/config"

	"github.com/andygeiss/cloud-native/utils/security"
)

func main() {
	cfg := config.Config{}
	srv := security.NewServer()

}
