package global

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/tal-tech/go-zero/core/logx"
	"sync/atomic"
	"time"
)

const (
	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"
)

type RedisFlowCountService struct {
	AppID       string        // 请求对应的appID 也可以是一种标识
	Interval    time.Duration // 确定多少时间间隔去统计一次流量总数
	QPS         int64         // 每秒请求量
	Unix        int64
	TickerCount int64
	TotalCount  int64 // 这个服务 当前时刻所在的这一天的流量
}

// 可以直接获取当前请求的qps 和 总流量

// NewRedisFlowCountService 为每一个服务创建一个流量统计器
func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		AppID:    appID,
		Interval: interval,
		QPS:      0,
		Unix:     0,
	}
	// 开一个协程 去根据时间间隔 去统计每个时间段的请求数
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Error(err)
			}
		}()
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			// 1.获取tickerCount 为 0
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) //获取数据 在这里读取value的值的同时，当前计算机中的任何CPU都不会进行其它的针对此值的读或写操作。
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据

			// 2.获取对应时间写入redis的key
			currentTime := time.Now()
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)

			// 3.写入每天和每小时的数据
			if err := RedisConfPipline(func(c redis.Conn) {
				c.Send("INCRBY", dayKey, tickerCount) // 将 key 所储存的值加上给定的增量值（increment）
				c.Send("EXPIRE", dayKey, 86400*2)     // 两天
				c.Send("INCRBY", hourKey, tickerCount)
				c.Send("EXPIRE", hourKey, 86400*2)
			}); err != nil {
				logx.Error("RedisConfPipline err", err)
				continue
			}

			// 4.获取这个时间段之前的总数
			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				logx.Error("reqCounter.GetDayData err", err)
				continue
			}

			// 5.获取到当前时间
			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue
			}

			tickerCount = totalCount - reqCounter.TotalCount // 获取到当前时间段的总数
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
			//fmt.Println("tickerCount", tickerCount, "totalCount", totalCount, "qps", reqCounter.QPS)
		}
	}()
	return reqCounter
}

// GetDayKey 提供一个key 这个key由这一天的日期和对应的服务ID组成   ( flow_day_count_20180808_service_id ) 最后加上服务名
func (o *RedisFlowCountService) GetDayKey(t time.Time) string {
	location, _ := time.LoadLocation("Asia/Chongqing")
	dayStr := t.In(location).Format("20060102")
	return fmt.Sprintf("%s_%s_%s", RedisFlowDayKey, dayStr, o.AppID)
}

// GetHourKey 提供一个小时的key ( flow_hour_count_2018080815_service_id )
func (o *RedisFlowCountService) GetHourKey(t time.Time) string {
	location, _ := time.LoadLocation("Asia/Chongqing")
	hourStr := t.In(location).Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", RedisFlowHourKey, hourStr, o.AppID)
}

// GetHourData 根据key获取对应小时的数据
func (o *RedisFlowCountService) GetHourData(hourTime time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", o.GetHourKey(hourTime)))
}

// GetDayData 根据key获取对应天的数据
func (o *RedisFlowCountService) GetDayData(dayTime time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", o.GetDayKey(dayTime)))
}

// Increase 原子增加 增加
func (o *RedisFlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		atomic.AddInt64(&o.TickerCount, 1)
	}()
}
