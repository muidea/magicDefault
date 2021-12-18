package dao

import "github.com/muidea/magicBatis/client"

type RemoteHub interface {
}

func New(clnt client.Client) RemoteHub {
	return &remoteHub{batisClient: clnt}
}

type remoteHub struct {
	batisClient client.Client
}
