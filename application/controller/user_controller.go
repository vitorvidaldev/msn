package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"github.com/vitorvidaldev/MSN/application/vo"
	"github.com/vitorvidaldev/MSN/domain/model"
	"github.com/vitorvidaldev/MSN/domain/repository"
	"github.com/vitorvidaldev/MSN/infra/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []vo.ReturnUserVO

	cur := repository.FindAll()
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cur, context.TODO())

	for cur.Next(context.TODO()) {
		var user model.User

		err := cur.Decode(&user)

		if err != nil {
			log.Fatal(err)
		}

		users = append(users, model.ToReturnVO(user))
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	id := extractID(r)
	err := repository.FindById(id, &user)

	if err != nil {
		util.GetError(err, w)
		return
	}

	err = json.NewEncoder(w).Encode(model.ToReturnVO(user))
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createUserVO vo.CreateUserVO
	_ = json.NewDecoder(r.Body).Decode(&createUserVO)

	user := model.FromCreateVO(createUserVO)

	result, err := repository.CreateUser(user)

	if err != nil {
		util.GetError(err, w)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updateUserVO vo.UpdateUserVO
	id := extractID(r)

	_ = json.NewDecoder(r.Body).Decode(&updateUserVO)

	user := model.FromUpdateVO(updateUserVO)

	update := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "username", Value: user.Username},
				primitive.E{Key: "email", Value: user.Email},
				primitive.E{Key: "updatedat", Value: time.Now()},
			},
		}}

	err := repository.UpdateUser(id, update, &user)

	if err != nil {
		util.GetError(err, w)
		return
	}

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := extractID(r)
	deleteResult, err := repository.DeleteUser(id)

	if err != nil {
		util.GetError(err, w)
		return
	}

	err = json.NewEncoder(w).Encode(deleteResult)
	if err != nil {
		log.Fatal(err)
	}
}

func extractID(r *http.Request) primitive.ObjectID {
	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
	}
	return id
}
