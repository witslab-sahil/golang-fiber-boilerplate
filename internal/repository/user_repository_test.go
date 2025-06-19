package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/witslab-sahil/fiber-boilerplate/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo UserRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repo = NewUserRepository(db)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM users")
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID)
}

func (suite *UserRepositoryTestSuite) TestGetByID() {
	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
	}
	suite.db.Create(user)

	found, err := suite.repo.GetByID(user.ID)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), found)
	assert.Equal(suite.T(), user.Email, found.Email)

	notFound, err := suite.repo.GetByID(99999)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), notFound)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
	}
	suite.db.Create(user)

	found, err := suite.repo.GetByEmail("test@example.com")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), found)
	assert.Equal(suite.T(), user.Username, found.Username)

	notFound, err := suite.repo.GetByEmail("notfound@example.com")
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), notFound)
}

func (suite *UserRepositoryTestSuite) TestGetByUsername() {
	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
	}
	suite.db.Create(user)

	found, err := suite.repo.GetByUsername("testuser")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), found)
	assert.Equal(suite.T(), user.Email, found.Email)

	notFound, err := suite.repo.GetByUsername("notfound")
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), notFound)
}

func (suite *UserRepositoryTestSuite) TestGetAll() {
	for i := 0; i < 5; i++ {
		user := &models.User{
			Email:    "test" + string(rune(i)) + "@example.com",
			Username: "testuser" + string(rune(i)),
			Password: "hashedpassword",
		}
		suite.db.Create(user)
	}

	users, total, err := suite.repo.GetAll(1, 3)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 3)
	assert.Equal(suite.T(), int64(5), total)

	users, total, err = suite.repo.GetAll(2, 3)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2)
	assert.Equal(suite.T(), int64(5), total)
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
	}
	suite.db.Create(user)

	user.FirstName = "Updated"
	user.LastName = "Name"
	err := suite.repo.Update(user)
	assert.NoError(suite.T(), err)

	updated, err := suite.repo.GetByID(user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated", updated.FirstName)
	assert.Equal(suite.T(), "Name", updated.LastName)
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	user := &models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
	}
	suite.db.Create(user)

	err := suite.repo.Delete(user.ID)
	assert.NoError(suite.T(), err)

	var count int64
	suite.db.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}