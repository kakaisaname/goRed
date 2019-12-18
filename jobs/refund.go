package jobs

import (
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/kakaisaname/infra"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"goRed/core/envelopes"
	"time"
)

//创建一个结构体，
type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker

	//互斥锁对象
	mutex *redsync.Mutex //分布式锁  redsync 开源组件
}

//初始化		**
func (r *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	d := ctx.Props().GetDurationDefault("jobs.refund.interval", time.Minute) //具体的时间间隔，在配置里获取		**
	//实例化 ticker
	r.ticker = time.NewTicker(d)
	maxIdle := ctx.Props().GetIntDefault("redis.maxIdle", 2)
	maxActive := ctx.Props().GetIntDefault("redis.maxActive", 5)
	timeout := ctx.Props().GetDurationDefault("redis.timeout", 20*time.Second)
	addr := ctx.Props().GetDefault("redis.addr", "127.0.0.1:6379")
	//
	////构建分布式锁   创建连接池
	pools := make([]redsync.Pool, 0)
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: timeout,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", addr) //创建tcp连接
		},
	}
	pools = append(pools, pool)
	rsync := redsync.New(pools)
	ip, err := utils.GetExternalIP() //这里获取当前的 Ipv4 地址
	if err != nil {
		ip = "127.0.0.1"
	}
	r.mutex = rsync.NewMutex("lock:RefundExpired",
		redsync.SetExpiry(50*time.Second),
		redsync.SetRetryDelay(3),
		redsync.SetGenValueFunc(func() (s string, e error) {
			now := time.Now()
			log.Infof("节点%s正在执行过期红包的退款任务", ip)
			return fmt.Sprintf("%d:%s", now.Unix(), ip), nil
		}),
	)
}

//在Start里用go协程来			for循环一直执行
func (r *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	go func() {
		for { //一直执行	**
			c := <-r.ticker.C //定时时间到期	**
			err := r.mutex.Lock()
			if err == nil { //拿到了锁 	**
				log.Debug("过期红包退款开始...", c)
				//红包过期退款的业务逻辑代码
				domain := envelopes.ExpiredEnvelopeDomain{}
				domain.Expired()
			} else {
				fmt.Println("c")
				log.Info("已经有节点在运行该任务了")
			}
			r.mutex.Unlock()

		}
	}()

}

func (r *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}

//这里把测试和debug留给同学们
