package global

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

func RedisConfPipline(pip ...func(c redis.Conn)) error {
	c, err := RedisConnFactory()
	if err != nil {
		return err
	}
	defer c.Close()
	for _, f := range pip {
		f(c)
	}
	c.Flush()
	return nil
}

func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	c, err := RedisConnFactory()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.Do(commandName, args...)
}

func RedisConnFactory() (redis.Conn, error) {
	c, err := redis.Dial(
		"tcp",
		"127.0.0.1:6379",
		redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond))
	if err != nil {
		return nil, err
	}

	//if _, err := c.Do("AUTH", ""); err != nil {
	//	c.Close()
	//	return nil, err
	//}

	//fmt.Println("连接redis成功")
	return c, nil
}
