package postgres

import "log"

func RunExample() {
	defer closeConnection()
	user := User{
		Username: "john.doe@gmail.com",
		Name:     "John",
		Surname:  "Doe",
		Role:     "Admin",
	}
	id, err := SaveUser(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user saved with id=%v", id)

	users, err := GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Printf("%v", user)
	}

	err = UpdateUser(User{
		ID:      id,
		Name:    "John2",
		Surname: "Doe",
		Role:    "User",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user updated with id=%v", id)

	users, err = GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Printf("%v", user)
	}

	err = DeleteUser(id)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user deleted with id=%v", id)
}
