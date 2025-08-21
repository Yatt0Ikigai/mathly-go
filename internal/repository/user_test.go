package repository_test

import (
	// "mathly/internal/models"
	"mathly/internal/models"
	"mathly/internal/repository"
	"time"

	// "time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", Ordered, func() {
	ID, _ := uuid.Parse("96c46ea4-3144-43fe-ab86-aab9f5c8de17")
	var (
		databases      repository.Databases
		userRepository repository.User
	)

	BeforeAll(func() {
		var err error
		databases, err = getDatabases()
		Expect(err).To(BeNil())
		repositories := repository.NewRepositories(databases)
		userRepository = repositories.User()
	})

	AfterAll(func() {
		databases.Close()
	})

	BeforeEach(func() {
		_, err := databases.DB().Query("DELETE from users")
		Expect(err).To(BeNil())
		_, err = databases.DB().Query(`
			INSERT INTO users (id, email, nickname, password_hash, created_at, updated_at) VALUES (
				'96c46ea4-3144-43fe-ab86-aab9f5c8de17', 'a@gmail.com', 'nickname', 'hash', NOW(), NOW()
			);
		`)
		Expect(err).To(BeNil())
	})

	Describe("GetByID", func() {
		It("should find a user by id in the db", func() {
			// given
			// when
			user, err := userRepository.GetByID(ID)
			// then
			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())
			Expect(user.ID).To(Equal(ID))
			Expect(user.Email).To(Equal("a@gmail.com"))
			Expect(user.Hash).To(Equal("hash"))
			Expect(user.Nickname).To(Equal("nickname"))
		})

		It("GetByEmail shouldn't find a user by id in the db", func() {
			// given
			// when
			user, err := userRepository.GetByEmail("invalidEmail")
			// then
			Expect(err).To(BeNil())
			Expect(user).To(BeNil())
		})
	})

	Describe("GetByEmail", func() {
		It("GetByEmail should find a user by id in the db", func() {
			// given
			// when
			user, err := userRepository.GetByEmail("a@gmail.com")
			// then
			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())
			Expect(user.ID).To(Equal(ID))
			Expect(user.Email).To(Equal("a@gmail.com"))
			Expect(user.Hash).To(Equal("hash"))
			Expect(user.Nickname).To(Equal("nickname"))
		})

		It("GetByEmail shouldn't find a user by id in the db", func() {
			// given
			// when
			user, err := userRepository.GetByEmail("invalidEmail")
			// then
			Expect(err).To(BeNil())
			Expect(user).To(BeNil())
		})
	})

	Describe("Insert", func() {
		It("Insert should find a create new user by id in the db", func() {
			// given
			u := models.User{
				ID:        uuid.Max,
				Email:     "some-email@gmail.com",
				Nickname:  "Yatt0",
				Hash:      "AAAAA",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			// when
			user, err := userRepository.Insert(&u)
			// then
			Expect(err).To(BeNil())
			dbUser, err := userRepository.GetByEmail(user.Email)
			Expect(err).To(BeNil())
			Expect(dbUser.ID).To(Equal(u.ID))
			Expect(dbUser.Email).To(Equal(u.Email))
			Expect(dbUser.Hash).To(Equal(u.Hash))
			Expect(dbUser.Nickname).To(Equal(u.Nickname))
			Expect(dbUser.Nickname).To(Equal(u.Nickname))
			Expect(dbUser.Nickname).To(Equal(u.Nickname))
		})
	})
})
