package sockets_test

import (
	"mathly/internal/log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepository(t *testing.T) {
	log.InitLogger()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sockets Suite")
}
