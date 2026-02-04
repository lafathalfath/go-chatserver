package cache

import (
	"encoding/json"
	"github.com/lafathalfath/go-chatserver/database"
	"time"
)

var conn = &database.DBConnection

func Scan(cacheKey string, dest any) error {
	val, err := conn.RDB.Get(conn.RDBCtx, cacheKey).Result()
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(val), &dest)
	return nil
}

func Set(cacheKey string, obj any, exp time.Duration) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	conn.RDB.Set(conn.RDBCtx, cacheKey, bytes, exp)
	return nil
}

func Del(cacheKey string) {
	_, err := conn.RDB.Get(conn.RDBCtx, cacheKey).Result()
	if err != nil {
		return
	}
	conn.RDB.Del(conn.RDBCtx, cacheKey)
}