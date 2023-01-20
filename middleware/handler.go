package middleware

import (
	"fmt"
	"net/http"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"crud/model"

	_ "github.com/lib/pq"
)

type response struct {
	Id int64 `json:"id"`
	Message string `json:"message"`
}


func CreateConnection() *sql.DB {
	
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error while loading env file")
	}

	conn, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB Connected Successfully")

	return conn
}

func CreateUser(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    res.Header().Set("Access-Control-Allow-Origin", "*")
    res.Header().Set("Access-Control-Allow-Methods", "POST")
    res.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user model.User

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Error while decoding %v", err)
	}

	userId := insertUser(user)

	response := response{
		Id: userId,
		Message: "Successfully added User",
	}

	json.NewEncoder(res).Encode(response)

}

func GetUser(res http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error while decoding request body %v", err)
	}
	
	user, err := getUser(int64(id))

	if err == sql.ErrNoRows {
		response := response {
			Id: int64(id),
			Message: "No data with given Id",
		}
		json.NewEncoder(res).Encode(response)
	} else if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Error while handling %v", err)
	} else {
		json.NewEncoder(res).Encode(user)
	}

}

func GetAllUsers(res http.ResponseWriter, req *http.Request) {
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Error While fetching data %v", err)
	}

	json.NewEncoder(res).Encode(users)
}

func UpdateUser(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert String to int, %v", err)
	}

	var user model.User

	err = json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Error while decoding request body, %v", err)
	}

	updatedRows := updateUserData(int64(id), user)

	var resp response
	if updatedRows == 0 {
		resp = response {
			Id : int64(id),
			Message: fmt.Sprintf("No row with given Id"),
		}
	} else {
		resp = response {
			Id : int64(id),
			Message: fmt.Sprintf("Total rows affected %d", updatedRows),
		}
	}

	json.NewEncoder(res).Encode(resp)
}


func DeleteUser(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error while decoding request body, %v", err)
	}

	deleteUser := deleteUserData(int64(id))

	var resp response
	if deleteUser == 0 {
		resp = response {
			Id : int64(id),
			Message: fmt.Sprintf("No row with given Id"),
		}
	} else {
		resp = response {
			Id : int64(id),
			Message: fmt.Sprintf("Total Rows deleted, %v", deleteUser),
		}
	}

	
	json.NewEncoder(res).Encode(resp)
}


func insertUser(user model.User) int64 {
	
	db := CreateConnection()

	defer db.Close()

	query := `INSERT INTO userdetails (name, email, contact) VALUES($1, $2, $3) RETURNING id`

	var id int64

	err := db.QueryRow(query, user.Name, user.Email, user.Contact).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute query %v", err)
	}

	fmt.Printf("Inserted a single record with %d", id)

	return id
}

func getUser(id int64) (model.User, error) {
	db := CreateConnection()

	defer db.Close()

	var user model.User

	query := `SELECT * FROM userdetails WHERE id=$1`

	row := db.QueryRow(query, id)
	
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Contact)

	// if err != nil {
	// 	log.Fatalf("Error while fetching row, %v", err)
	// }
	// if err == sql.ErrNoRows {
	// 	fmt.Println("No row with given id")
	// 	return user, nil
	// }

	// switch err {
    // case sql.ErrNoRows:
    //     fmt.Println("No rows were returned!")
    //     return user, nil
    // case nil:
    //     return user, nil
    // default:
    //     log.Fatalf("Unable to scan the row. %v", err)
    // }
	
	return user, err
}

func getAllUsers() ([]model.User, error) {
	db := CreateConnection()

	defer db.Close()

	query := `SELECT * FROM userdetails`

	var users []model.User

	rows, err := db.Query(query)

	if err != nil {
		log.Fatalf("Error while querying. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Contact)

		if err != nil {
			log.Fatalf("Error while fetching data. %v", err)
		}

		users = append(users, user)
	}	

	return users, err
	
}

func updateUserData(id int64, user model.User) int64 {
	db := CreateConnection()

	defer db.Close()

	query := `UPDATE userdetails SET name=$2, email=$3, contact=$4 WHERE id=$1`

	res, err := db.Exec(query, id, user.Name, user.Email, user.Contact)

	if err != nil {
		log.Fatalf("Unable to Update data %v", err);
	}
	
	updatedRows, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Unable to determine rows Affected, %v", err)
	}

	return updatedRows

}

func deleteUserData(id int64) int64 {
	db := CreateConnection()

	defer db.Close()

	query := `DELETE FROM userdetails WHERE id=$1`

	res, err := db.Exec(query, id)

	if err != nil {
		log.Fatalf("Unable to delete data, %v", err)
	}

	deletedRows, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Unable to determine deleted data, %v", err)
	}

	return deletedRows
}