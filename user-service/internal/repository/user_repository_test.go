package repository_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/NeGat1FF/e-commerce/user-service/internal/models"
	"github.com/NeGat1FF/e-commerce/user-service/internal/repository"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ps "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestContainer(ctx context.Context) (*gorm.DB, error) {
	pgCon, err := ps.Run(ctx,
		"postgres:latest",
		ps.WithDatabase("test"),
		ps.WithUsername("test"),
		ps.WithPassword("test"),
		ps.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	connString, err := pgCon.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("postgres://test:test@%s/test?sslmode=disable", connString)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateUser(t *testing.T) {

	ctx := context.Background()
	db, err := setupTestContainer(ctx)
	require.NoError(t, err)

	repo := repository.NewUserRepository(db)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Create user",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := repo.CreateUser(context.Background(), &user)
				assert.NoError(t, err)
			},
		},
		{
			name: "Create user with same email",
			testFunc: func(t *testing.T) {
				firstUser := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := repo.CreateUser(context.Background(), &firstUser)
				assert.NoError(t, err)

				secondUser := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    firstUser.Email,
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err = repo.CreateUser(context.Background(), &secondUser)
				assert.Error(t, err)
				t.Log(err)
			},
		},
		{
			name: "Create user with field too long",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     strings.Repeat("a", 256),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := repo.CreateUser(context.Background(), &user)
				assert.Error(t, err)
			},
		},
		{
			name: "Create user with empty field",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := repo.CreateUser(context.Background(), &user)
				assert.Error(t, err)
			},
		},
		{
			name: "Create user with nil object",
			testFunc: func(t *testing.T) {
				err := repo.CreateUser(context.Background(), nil)
				assert.Error(t, err)
			},
		},
		{
			name: "Create user with cancel context",
			testFunc: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := repo.CreateUser(ctx, &user)
				assert.Error(t, err)
			},
		},
		{
			name: "Create users in parallel",
			testFunc: func(t *testing.T) {
				users := make([]models.User, 10)
				for i := range users {
					users[i] = models.User{
						ID:       uuid.New(),
						Name:     gofakeit.FirstName(),
						Surname:  gofakeit.LastName(),
						Email:    gofakeit.Email(),
						Password: gofakeit.Password(true, true, true, true, true, 10),
						Phone:    gofakeit.Phone(),
					}
				}

				wg := sync.WaitGroup{}
				wg.Add(len(users))
				errs := make(chan error, len(users))
				for i := range users {
					go func(i int) {
						defer wg.Done()
						errs <- repo.CreateUser(context.Background(), &users[i])
					}(i)
				}

				for range users {
					err := <-errs
					assert.NoError(t, err)
				}

				wg.Wait()
			},
		},
	}

	for _, tc := range testCases {
		err = db.AutoMigrate(&models.User{})
		require.NoError(t, err)

		t.Run(tc.name, tc.testFunc)

		err = db.Migrator().DropTable(&models.User{})
		require.NoError(t, err)
	}

}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	db, err := setupTestContainer(ctx)
	require.NoError(t, err)

	repo := repository.NewUserRepository(db)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Get user by ID",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := db.Create(&user).Error
				require.NoError(t, err)

				userFromDB, err := repo.GetUserByID(context.Background(), user.ID.String())

				assert.Equal(t, user.Name, userFromDB.Name)
				assert.NoError(t, err)
			},
		},
		{
			name: "Get user by ID with invalid ID",
			testFunc: func(t *testing.T) {
				user, err := repo.GetUserByID(context.Background(), "invalid")

				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		err = db.AutoMigrate(&models.User{})
		require.NoError(t, err)

		t.Run(tc.name, tc.testFunc)

		err = db.Migrator().DropTable(&models.User{})
		require.NoError(t, err)
	}

}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	db, err := setupTestContainer(ctx)
	require.NoError(t, err)

	repo := repository.NewUserRepository(db)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Get user by email",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := db.Create(&user).Error
				require.NoError(t, err)

				userFromDB, err := repo.GetUserByEmail(context.Background(), user.Email)

				assert.Equal(t, user.Name, userFromDB.Name)
				assert.NoError(t, err)
			},
		},
		{
			name: "Get user by ID with invalid email",
			testFunc: func(t *testing.T) {
				user, err := repo.GetUserByEmail(context.Background(), "invalid")

				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		err = db.AutoMigrate(&models.User{})
		require.NoError(t, err)

		t.Run(tc.name, tc.testFunc)

		err = db.Migrator().DropTable(&models.User{})
		require.NoError(t, err)
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	db, err := setupTestContainer(ctx)
	require.NoError(t, err)

	repo := repository.NewUserRepository(db)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Update user",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := db.Create(&user).Error
				require.NoError(t, err)

				updatedUser := models.User{
					ID:   user.ID,
					Name: gofakeit.FirstName(),
				}

				err = repo.UpdateUser(context.Background(), &updatedUser)
				require.NoError(t, err)

				var userFromDB models.User
				db.First(&userFromDB, "id = ?", user.ID)

				assert.Equal(t, updatedUser.Name, userFromDB.Name)
				assert.Equal(t, user.Surname, userFromDB.Surname)
			},
		},
		{
			name: "Update user with invalid ID",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := db.Create(&user).Error
				require.NoError(t, err)

				user.ID = uuid.New()

				err = repo.UpdateUser(context.Background(), &user)
				require.Error(t, err)
			},
		},
		{
			name: "Update user with field too long",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := db.Create(&user).Error
				require.NoError(t, err)

				updatedUser := models.User{
					ID:   user.ID,
					Name: strings.Repeat("a", 256),
				}

				err = repo.UpdateUser(context.Background(), &updatedUser)
				require.Error(t, err)
			},
		},
		{
			name: "Update user with nil object",
			testFunc: func(t *testing.T) {
				err := repo.UpdateUser(context.Background(), nil)
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		err = db.AutoMigrate(&models.User{})
		require.NoError(t, err)

		t.Run(tc.name, tc.testFunc)

		err = db.Migrator().DropTable(&models.User{})
		require.NoError(t, err)
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	db, err := setupTestContainer(ctx)
	require.NoError(t, err)

	repo := repository.NewUserRepository(db)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Delete user",
			testFunc: func(t *testing.T) {
				user := models.User{
					ID:       uuid.New(),
					Name:     gofakeit.FirstName(),
					Surname:  gofakeit.LastName(),
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, true, 10),
					Phone:    gofakeit.Phone(),
				}

				err := db.Create(&user).Error
				require.NoError(t, err)

				err = repo.DeleteUser(context.Background(), user.ID.String())
				require.NoError(t, err)

				var userFromDB models.User
				db.First(&userFromDB, "id = ?", user.ID)
				assert.Equal(t, models.User{}, userFromDB)
			},
		},
		{
			name: "Delete user with invalid ID",
			testFunc: func(t *testing.T) {
				err := repo.DeleteUser(context.Background(), "invalid")
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		err = db.AutoMigrate(&models.User{})
		require.NoError(t, err)

		t.Run(tc.name, tc.testFunc)

		err = db.Migrator().DropTable(&models.User{})
		require.NoError(t, err)
	}
}
