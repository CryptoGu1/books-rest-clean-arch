package log_grpc

import (
	"context"
	"fmt"

	audit "github.com/CryptoGu1/books-grpc-log/pkg/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn        *grpc.ClientConn
	auditClient audit.AuditServiceClient
}

func NewClient(host string, port int) (*Client, error) {
	var conn *grpc.ClientConn

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:        conn,
		auditClient: audit.NewAuditServiceClient(conn),
	}, nil

}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) SendLogRequest(ctx context.Context, req audit.LogItem) error {
	action, err := audit.ToPbAction(req.Action)
	if err != nil {
		return err
	}

	entity, err := audit.ToPbEntity(req.Entity)
	if err != nil {
		return err
	}

	_, err = c.auditClient.Log(ctx, &audit.LogRequest{
		Action:    action,
		Entity:    entity,
		IntityId:  req.EntityID,
		Timestamp: timestamppb.New(req.Timestamp),
	})

	return err
}
