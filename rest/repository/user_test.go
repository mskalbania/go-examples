package repository

import (
	"context"
	"fmt"
	"github.com/Shopify/toxiproxy/client"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"go-examples/rest/config"
	"go-examples/rest/database"
	"go-examples/rest/model"
	"testing"
	"time"
)

var testUser = model.PostUser{
	Email: "test@example.org",
}
var timeout = 50 * time.Millisecond

type UserSuite struct {
	suite.Suite
	userRepository     *UserRepository
	closeDb            func()
	database           database.Database
	postgresContainer  *postgres.PostgresContainer
	toxiproxyContainer testcontainers.Container
	postgresProxy      *toxiproxy.Proxy
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) SetupSuite() {
	net, err := network.New(context.Background())
	if err != nil {
		suite.T().Fatal(err)
	}
	if toxiCt, proxy, err := SetupToxiproxyContainer(net); err != nil {
		suite.T().Fatal(err)
	} else {
		suite.toxiproxyContainer = toxiCt
		suite.postgresProxy = proxy
	}
	if pgCt, err := SetupPostgresContainer(net); err != nil {
		suite.T().Fatal(err)
	} else {
		suite.postgresContainer = pgCt
	}
	conf, err := CreateDBProperties(suite.toxiproxyContainer)
	if err != nil {
		suite.T().Fatal(err)
	}
	db, cancel, err := database.NewPostgresDatabase(&config.AppConfig{DB: conf})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.database = db
	suite.closeDb = cancel
	suite.userRepository = NewUserRepository(db, &conf)
}

func SetupPostgresContainer(net *testcontainers.DockerNetwork) (*postgres.PostgresContainer, error) {
	pgCt, err := postgres.Run(context.Background(),
		"postgres:16-alpine",
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithDatabase("postgres"),
		postgres.WithInitScripts("../docker/schema.sql"),
		postgres.BasicWaitStrategies(),
		network.WithNetwork([]string{"postgres"}, net),
	)
	if err != nil {
		return nil, err
	}
	return pgCt, nil
}

func SetupToxiproxyContainer(net *testcontainers.DockerNetwork) (testcontainers.Container, *toxiproxy.Proxy, error) {
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "shopify/toxiproxy",
			Networks:     []string{net.Name},
			ExposedPorts: []string{"8474/tcp", "8666/tcp"}, //control port / proxy port
			WaitingFor:   wait.ForHTTP("/version").WithPort("8474/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}
	port, err := container.MappedPort(context.Background(), "8474")
	if err != nil {
		return nil, nil, err
	}
	client := toxiproxy.NewClient(fmt.Sprintf("http://localhost:%d", port.Int()))
	//listen on 8666 exposed port from toxi proxy container
	//forward to 5432 postgres internal port, host - network alias of postgres container
	proxy, err := client.CreateProxy("postgres", "0.0.0.0:8666", "postgres:5432")
	if err != nil {
		return nil, nil, err
	}
	return container, proxy, nil
}

func CreateDBProperties(toxiCt testcontainers.Container) (config.DBConfig, error) {
	host, err := toxiCt.Host(context.Background())
	if err != nil {
		return config.DBConfig{}, err
	}
	port, err := toxiCt.MappedPort(context.Background(), "8666")
	if err != nil {
		return config.DBConfig{}, err
	}
	return config.DBConfig{
		User:     "postgres",
		Password: "postgres",
		Host:     host,
		Port:     port.Int(),
		Database: "postgres",
		Timeout:  timeout,
		PoolMin:  1,
		PoolMax:  1,
	}, nil
}

func (suite *UserSuite) TearDownSuite() {
	suite.closeDb()
	_ = suite.postgresContainer.Terminate(context.Background())
	_ = suite.toxiproxyContainer.Terminate(context.Background())
}

func (suite *UserSuite) TearDownTest() {
	_, err := suite.userRepository.database.Exec(context.Background(), "TRUNCATE public.user")
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *UserSuite) TestSaveUser() {
	//when
	saved, err := suite.userRepository.Save(context.Background(), &testUser)

	//then
	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), saved.ID)
	require.Equal(suite.T(), testUser.Email, saved.Email)

	//and
	get, err := suite.userRepository.GetUserById(context.Background(), saved.ID)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), saved.ID, get.ID)
	require.Equal(suite.T(), saved.Email, get.Email)
}

func (suite *UserSuite) TestGetAllUsers() {
	//given
	saved1, _ := suite.userRepository.Save(context.Background(), &testUser)
	saved2, _ := suite.userRepository.Save(context.Background(), &testUser)

	//when
	users, err := suite.userRepository.GetAllUsers(context.Background())

	//then
	require.NoError(suite.T(), err)
	require.Len(suite.T(), users, 2)
	require.Contains(suite.T(), users, saved1)
	require.Contains(suite.T(), users, saved2)
}

func (suite *UserSuite) TestGetUserById() {
	//given
	saved, _ := suite.userRepository.Save(context.Background(), &testUser)

	//when
	get, err := suite.userRepository.GetUserById(context.Background(), saved.ID)

	//then
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), saved.ID, get.ID)
	require.Equal(suite.T(), saved.Email, get.Email)
}

func (suite *UserSuite) TestExists() {
	//given
	saved, _ := suite.userRepository.Save(context.Background(), &testUser)

	//when
	exists, err := suite.userRepository.Exists(context.Background(), saved.ID)

	//then
	require.NoError(suite.T(), err)
	require.True(suite.T(), exists)
}

func (suite *UserSuite) TestDelete() {
	//given
	saved, _ := suite.userRepository.Save(context.Background(), &testUser)

	//when
	err := suite.userRepository.Delete(context.Background(), saved.ID)

	//then
	require.NoError(suite.T(), err)

	//and
	exists, err := suite.userRepository.Exists(context.Background(), saved.ID)
	require.NoError(suite.T(), err)
	require.False(suite.T(), exists)
}

func (suite *UserSuite) TestUpdate() {
	//given
	saved, _ := suite.userRepository.Save(context.Background(), &testUser)
	updateRq := model.PostUser{
		Email: "new@gmail.com",
	}

	//when
	updated, err := suite.userRepository.Update(context.Background(), saved.ID, &updateRq)

	//then
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), saved.ID, updated.ID)
	require.Equal(suite.T(), updateRq.Email, updated.Email)
}

func (suite *UserSuite) TestTimeout() {
	//given
	_, err := suite.postgresProxy.AddToxic("postgres", "latency", "downstream", 1.0,
		toxiproxy.Attributes{"latency": timeout.Milliseconds() + 100},
	)
	if err != nil {
		suite.T().Fatal(err)
	}

	cases := []struct {
		operationName string
		operationF    func() error
	}{
		{
			operationName: "Save",
			operationF: func() error {
				_, err := suite.userRepository.Save(context.Background(), &testUser)
				return err
			},
		},
		{
			operationName: "GetAllUsers",
			operationF: func() error {
				_, err := suite.userRepository.GetAllUsers(context.Background())
				return err
			},
		},
		{
			operationName: "GetUserById",
			operationF: func() error {
				_, err := suite.userRepository.GetUserById(context.Background(), "1")
				return err
			},
		},
		{
			operationName: "Delete",
			operationF: func() error {
				return suite.userRepository.Delete(context.Background(), "1")
			},
		},
		{
			operationName: "Exists",
			operationF: func() error {
				_, err := suite.userRepository.Exists(context.Background(), "1")
				return err
			},
		},
		{
			operationName: "Update",
			operationF: func() error {
				_, err := suite.userRepository.Update(context.Background(), "1", &testUser)
				return err
			},
		},
	}

	//when
	for _, c := range cases {
		err := c.operationF()

		//then
		require.Error(suite.T(), err)
		require.Contains(suite.T(), err.Error(), "context deadline exceeded", "%s should have timed out", c.operationName)
	}
	_ = suite.postgresProxy.RemoveToxic("postgres")
}
