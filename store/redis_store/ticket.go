package redis_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

type redisTicket struct {
	client *redis.Client
}

func (rt *redisTicket) TicketIsNotExistErr(err error) bool {
	return model.TicketIsNotExistErr(err)
}

func (rt *redisTicket) TicketLoad(id string) (ticket *model.Ticket, err error) {

	if id == "" {
		return nil, model.ErrTicketNotExist
	}

	res, err := rt.client.Get(id).Result()
	if err != nil {
		if err == redis.Nil {
			err = model.ErrTicketNotExist
		}
		return nil, err
	}
	ticket = &model.Ticket{}
	if err = msgpack.Unmarshal([]byte(res), ticket); err != nil {
		return nil, err
	}
	return
}

func (rt *redisTicket) TicketCreate(ticket *model.Ticket) error {
	if res, err := rt.client.Exists(ticket.Id).Result(); err != nil {
		return err
	} else if res != 0 {
		return model.ErrTicketExisted
	}

	b, err := msgpack.Marshal(ticket)
	if err != nil {
		return err
	}
	return rt.client.Set(ticket.Id, b, ticket.ExpiredAt.Sub(ticket.CreatedAt)).Err()
}

func (rt *redisTicket) TicketDelete(id string) error {
	return rt.client.Del(id).Err()
}

func NewRedisTicket(client *redis.Client) model.TicketStore {
	return &redisTicket{client: client}
}
