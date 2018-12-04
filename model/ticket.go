package model

import (
	"errors"
	"time"
)

// 登录凭证
type Ticket struct {
	Id        string
	UserId    int64
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TicketStore interface {
	TicketLoad(id string) (*Ticket, error)
	TicketCreate(*Ticket) error
	TicketDelete(id string) error
	TicketIsNotExistErr(error) bool
}

type TicketService interface {
	TicketIsValid(ticketId string) (isValid bool, userId int64, err error)
	// 生成 ticket
	TicketGen(userId int64) (*Ticket, error)
	TicketTTL() time.Duration
	TicketDestroy(ticketId string) error
}

var (
	ErrTicketNotExist = errors.New("ticket not exist")
	ErrTicketExisted  = errors.New("ticket 已经存在")
)

func TicketIsNotExistErr(err error) bool {
	return err == ErrTicketNotExist
}
