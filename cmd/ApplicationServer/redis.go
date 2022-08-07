package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/SiddeshSambasivam/shillings/pkg/models"
)

func (env *DataEnv) SetUserProfCache(
	ctx context.Context, usr models.User, claims models.Claims) error {

	marUsr, _ := json.Marshal(usr)
	var key string = strconv.FormatInt(int64(claims.User_id), 10)
	ttl := time.Until(time.Unix(claims.StandardClaims.ExpiresAt, 0))

	err := env.Redis.Set(ctx, key, marUsr, ttl).Err()
	if err != nil {
		log.Println("Error setting user profile in cache: ", err)
		return err
	}

	return nil
}

func (env *DataEnv) GetUserProfCache(ctx context.Context, claims models.Claims) (models.User, error) {

	usr := models.User{}
	var key string = strconv.FormatInt(int64(claims.User_id), 10)

	strcmd := env.Redis.Get(ctx, key)
	redisUsr, err := strcmd.Result()

	if redisUsr == "" {
		return usr, errors.New("user not found")
	} else if err != nil {
		return usr, err
	}

	err = json.Unmarshal([]byte(redisUsr), &usr)
	if err != nil {
		return usr, err
	}

	return usr, nil
}

func (env *DataEnv) DelUserProfCache(ctx context.Context, claims models.Claims) {
	env.Redis.Del(ctx, string(claims.User_id))
}
