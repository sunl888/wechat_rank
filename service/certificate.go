package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/gin-gonic/gin"
)

type certificateService struct {
	cs model.CertificateStore
}

func (cSvc *certificateService) CertificateUpdate(oldAccount, newAccount string, certificateType model.CertificateType) error {
	return cSvc.cs.CertificateUpdate(oldAccount, newAccount, certificateType)
}

func CertificateUpdate(ctx *gin.Context, oldAccount, newAccount string, certificateType model.CertificateType) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.CertificateUpdate(oldAccount, newAccount, certificateType)
	}
	return nil
}

func NewCertificateService(cs model.CertificateStore) model.CertificateService {
	return &certificateService{cs: cs}
}
