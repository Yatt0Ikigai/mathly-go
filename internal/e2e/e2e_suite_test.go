package e2e_test

import (
	"mathly/internal/log"
	"mathly/internal/repository"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepository(t *testing.T) {
	if os.Getenv("TEST_TYPE") == "e2e" {
		log.InitLogger()
		RegisterFailHandler(Fail)
		RunSpecs(t, "Repository Suite")
	}
}

func getDatabases() (repository.Databases, error) {
	return repository.NewDatabases(&repository.DatabasesConfig{
		Redis: repository.RedisConfig{},
		SQL: repository.SQLConfig{
			Host:     "localhost",
			Port:     5432,
			DB:       "mydatabase",
			User:     "myuser",
			Password: "mypassword",
			Schema:   "mathly",
		},
	})
}
