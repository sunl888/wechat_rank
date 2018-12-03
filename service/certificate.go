package service

import "code.aliyun.com/zmdev/wechat_rank/model"

type certificateService struct {
	cs model.CertificateStore
}

func (cSvc *certificateService) CertificateUpdate(oldAccount, newAccount string, certificateType model.CertificateType) error {
	return cSvc.cs.CertificateUpdate(oldAccount, newAccount, certificateType)
}

func NewCertificateService(cs model.CertificateStore) model.CertificateService {
	return &certificateService{cs: cs}
}
