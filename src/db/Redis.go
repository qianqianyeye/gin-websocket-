package db

import (
	"github.com/go-redis/redis"
	"fmt"
)

var CfRedis []*redis.Client

func InitRedis(flag bool){
	if len(CfRedis)>0 {
		fmt.Println(flag)
		CfRedis=append(CfRedis[:0],CfRedis[len(CfRedis):]...)
	}
	var s []*redis.Client
	if flag {
		adPush := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "4%CPpOoUPML0&SMa", // no password set
				DB:       3,                  // use default DB)
			})
		clawGlad := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "4%CPpOoUPML0&SMa", // no password set
				DB:       2,                  // use default DB)
			})
		s = append(s, adPush)
		s = append(s, clawGlad)
	}else {
		adPush := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "4%CPpOoUPML0&SMa", // no password set
				DB:       5,                  // use default DB)
			})
		clawGlad := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "4%CPpOoUPML0&SMa", // no password set
				DB:       10,                  // use default DB)
			})
		s = append(s, adPush)
		s = append(s, clawGlad)
	}
	CfRedis =s
}

func GetClientRedis() []*redis.Client {
	return CfRedis
}
