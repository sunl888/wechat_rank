package model

import (
	"errors"
)

type CertificateType uint8

const (
	UserName CertificateType = iota
	PhoneNum
	Email
)

type Certificate struct {
	Id      int64
	UserId  int64
	Account string `gorm:"not null;unique"`
	Type    CertificateType
}

type CertificateStore interface {
	CertificateExist(account string) (bool, error)
	CertificateLoadByAccount(account string) (*Certificate, error)
	CertificateIsNotExistErr(error) bool
	CertificateCreate(*Certificate) error
	CertificateUpdate(oldAccount, newAccount string, certificateType CertificateType) error
}

var ErrCertificateNotExist = errors.New("certificate not exist")

func CertificateIsNotExistErr(err error) bool {
	return err == ErrCertificateNotExist
}

type CertificateService interface {
	CertificateUpdate(oldAccount, newAccount string, certificateType CertificateType) error
}
