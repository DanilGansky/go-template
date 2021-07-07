package repository_test

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/littlefut/go-template/pkg/errors"

	"github.com/littlefut/go-template/internal/repository"
	"github.com/littlefut/go-template/internal/user"
	"github.com/littlefut/go-template/pkg/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

const defaultTimeout = time.Second * 3

var _ = Describe("User Repository", func() {
	type testData struct {
		id                  int
		username, password  string
		joinedAt, lastLogin time.Time
	}

	var (
		conn    *gorm.DB
		repo    user.Repository
		compose *testcontainers.LocalDockerCompose

		dockerComposeLocation = []string{"../../deploy/test/docker-compose.yml"}
		dsn                   = "host=localhost port=5433 user=postgres password=postgres dbname=service_db_test sslmode=disable"
	)

	BeforeSuite(func() {
		compose = testcontainers.NewLocalDockerCompose(dockerComposeLocation, "go-template-test")
		execErr := compose.WithCommand([]string{"up", "-d", "--remove-orphans"}).Invoke()
		Expect(execErr.Error).NotTo(HaveOccurred())

		time.Sleep(defaultTimeout)
		conn = db.Get(dsn)
		repo = repository.NewUserRepository(conn)
	})

	Context("FindByID", func() {
		var td testData

		BeforeEach(func() {
			td = testData{
				id:        1,
				username:  "test",
				password:  "$2a$12$GYOWNlR5V7Qqosfo.gomXey7H/WrGWNCV3MTlHD1WHFVYbx5EHMNS",
				joinedAt:  time.Now(),
				lastLogin: time.Now(),
			}

			Expect(conn.Exec(`INSERT INTO users (id, username, password, joined_at, last_login) VALUES (?, ?, ?, ?, ?)`,
				td.id, td.username, td.password, td.joinedAt, td.lastLogin).Error).NotTo(HaveOccurred())
		})

		When("user exists", func() {
			It("should return user", func() {
				ctx := context.Background()
				user, err := repo.FindByID(ctx, td.id)
				Expect(err).NotTo(HaveOccurred())
				Expect(user.ID).To(Equal(td.id))
				Expect(user.Username).To(Equal(td.username))
				Expect(user.Password).To(Equal(td.password))
			})
		})

		When("user does not exists", func() {
			It("should return nil and error NotFound", func() {
				ctx := context.Background()
				user, err := repo.FindByID(ctx, 0)
				Expect(user).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.(errors.Error).ErrorCode()).To(Equal(errors.NotFoundError))
			})
		})

		AfterEach(func() {
			Expect(conn.Exec(`DELETE FROM users WHERE username = ?`, td.username).Error).NotTo(HaveOccurred())
		})
	})

	AfterSuite(func() {
		db, _ := conn.DB()
		Expect(db.Close()).NotTo(HaveOccurred())
		execErr := compose.Down()
		Expect(execErr.Error).NotTo(HaveOccurred())
	})
})
