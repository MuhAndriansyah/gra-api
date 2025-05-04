package paymentgateway

import (
	"backend-layout/internal/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	SnapClient *snap.Client
	Coreapi    *coreapi.Client
}

func InitMidtrans(conf config.MidtransConfig) *MidtransClient {

	var s = snap.Client{}
	var c = coreapi.Client{}

	env := midtrans.Sandbox

	if conf.Mode == "production" {
		env = midtrans.Production
	}

	s.New(conf.ServerKey, env)
	c.New(conf.ServerKey, env)

	return &MidtransClient{
		SnapClient: &s,
		Coreapi:    &c,
	}
}
