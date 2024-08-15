package postgres

import (
	"context"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ct, err := postgres.Run(context.Background(),
		"postgres:16-alpine",
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts("docker/init.sql"),
		postgres.BasicWaitStrategies(), //to not try to connect before postgres is ready
	)

	if err != nil {
		log.Fatalf("error pooling postgres container: %v", err)
	}

	if err := ct.Start(context.Background()); err != nil {
		log.Fatalf("error starting postgres container: %v", err)
	}

	//get random host port allocated pointing to container's 5432 port
	if port, err := ct.MappedPort(context.Background(), "5432"); err != nil {
		log.Fatalf("error getting mapped port: %v", err)
	} else {
		openConnection("postgres", "postgres", "localhost", port.Int(), "postgres")
	}

	code := m.Run()

	closeConnection()
	if err := ct.Terminate(context.Background()); err != nil {
		log.Printf("error terminating postgres container: %v", err)
	}
	os.Exit(code)
}

// This is one "bulk" test to prove CRUD works.
// Not ideal approach, but works for demonstration purposes.
func TestSaveGetUpdateDelete(t *testing.T) {
	t.Log("Given user")
	{
		user := User{
			Username: "test@gmail.com",
			Name:     "testName",
			Surname:  "testSurname",
			Role:     "Admin",
		}
		t.Log("\tWhen user saved")
		{
			id, err := SaveUser(user)
			if err != nil {
				t.Fatalf("\t\tThen error not expected, actual: %v", err)
			} else {
				t.Log("\t\tThen user successfully saved")
			}
			if id == "" {
				t.Fatalf("\t\tThen user id expected, actual: %v", id)
			} else {
				t.Logf("\t\tThen user id is assigned: %v", id)
				user.ID = id
			}
		}
		t.Log("\tWhen all users retrieved")
		{
			users, err := GetAllUsers()
			if err != nil {
				t.Fatalf("\t\tThen error not expected, actual: %v", err)
			} else {
				t.Log("\t\tThen users successfully retrieved")
			}
			if len(users) != 1 {
				t.Fatalf("\t\tThen one user expected, actual: %v", len(users))
			} else {
				t.Log("\t\tThen one user retrieved")
			}
			if userDataMatch(*users[0], user) {
				t.Fatalf("\t\tThen user expected %v, actual: %v", user, users[0])
			} else {
				t.Log("\t\tThen retrieved user data matches expected")
			}
		}
		newUserData := User{
			ID:       user.ID,
			Username: user.Username,
			Name:     "newName",
			Surname:  "newSurname",
			Role:     "User",
		}
		t.Log("\tWhen user updated")
		{
			err := UpdateUser(newUserData)
			if err != nil {
				t.Fatalf("\t\tThen error not expected, actual: %v", err)
			} else {
				t.Log("\t\tThen user successfully updated")
			}
			users, _ := GetAllUsers()
			if userDataMatch(*users[0], newUserData) {
				t.Fatalf("\t\tThen new user data expected %v, actual: %v", newUserData, users[0])
			} else {
				t.Log("\t\tThen new retrieved user data matches expected")
			}
		}
		t.Log("\tWhen user deleted")
		{
			err := DeleteUser(user.ID)
			if err != nil {
				t.Fatalf("\t\tThen error not expected, actual: %v", err)
			} else {
				t.Log("\t\tThen user successfully deleted")
			}
			users, _ := GetAllUsers()
			if len(users) == 0 {
				t.Log("\t\tThen no users retrieved")
			} else {
				t.Fatalf("\t\tThen users not expected but present")
			}
		}
	}
}

func userDataMatch(actual User, expected User) bool {
	return actual.Username != expected.Username ||
		actual.Name != expected.Name ||
		actual.Surname != expected.Surname ||
		actual.Role != expected.Role
}
