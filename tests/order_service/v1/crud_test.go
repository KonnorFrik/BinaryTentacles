package order_service_v1_test

import (
	"context"
	"io"
	"testing"
	"time"

	client "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	marketIdValid    = "5d6f8857-fafe-432c-8380-2b340ec03bb7"
	orderServiceAddr = "0.0.0.0:8888"
)

var (
	orderService client.OrderServiceClient
	baseCtx      = context.Background()
	orderId      string
	userID       string
)

func init() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(orderServiceAddr, opts...)

	if err != nil {
		panic(err)
	}

	orderService = client.NewOrderServiceClient(conn)
}

func init() {
	id, err := uuid.NewV7()

	if err != nil {
		panic(err)
	}

	userID = id.String()
}

func TestCRUD(t *testing.T) {
	cases := []struct {
		name string
		f    func(*testing.T)
	}{
		{
			name: "Create order with valid market",
			f: func(t *testing.T) {
				req := client.CreateRequest{
					UserId:    userID,
					MarketId:  marketIdValid,
					OrderType: client.OrderType_ORDER_TYPE_T1,
					Price:     123,
					Quantity:  1,
				}
				resp, err := orderService.Create(baseCtx, &req)

				if err != nil {
					t.Errorf("Got = %q\n", err)
				}

				if resp.GetOrderStatus() != client.OrderStatus_ORDER_STATUS_CREATED {
					t.Errorf("Got = %d, Want = %d\n", resp.GetOrderStatus(), client.OrderStatus_ORDER_STATUS_CREATED)
				}

				orderId = resp.GetOrderId()
			},
		},

		{
			name: "Get status",
			f: func(t *testing.T) {
				req := client.OrderStatusRequest{
					OrderId: orderId,
					UserId:  userID,
				}
				resp, err := orderService.OrderStatus(baseCtx, &req)

				if err != nil {
					t.Fatalf("Got = %q\n", err)
				}

				if resp.GetStatus() != client.OrderStatus_ORDER_STATUS_CREATED {
					t.Fatalf("Got = %d, Want = %d\n", resp.GetStatus(), client.OrderStatus_ORDER_STATUS_CREATED)
				}
			},
		},

		{
			name: "Create order with invalid market",
			f: func(t *testing.T) {
				req := client.CreateRequest{
					UserId:    userID,
					MarketId:  "",
					OrderType: client.OrderType_ORDER_TYPE_T1,
					Price:     123,
					Quantity:  1,
				}
				_, err := orderService.Create(baseCtx, &req)

				if err == nil {
					t.Fatalf("Got nil error")
				}

				stat, ok := status.FromError(err)

				if !ok {
					t.Fatalf("Error on convert status from error: %q\n", err)
				}

				if stat.Code() != codes.FailedPrecondition {
					t.Fatalf("Got = %d, Want = %d\n", stat.Code(), codes.FailedPrecondition)
				}
			},
		},

		{
			name: "Create order with invalid market",
			f: func(t *testing.T) {
				req := client.CreateRequest{
					UserId:    userID,
					MarketId:  "1234",
					OrderType: client.OrderType_ORDER_TYPE_T1,
					Price:     123,
					Quantity:  1,
				}
				_, err := orderService.Create(baseCtx, &req)

				if err == nil {
					t.Fatalf("Got nil error")
				}

				stat, ok := status.FromError(err)

				if !ok {
					t.Fatalf("Error on convert status from error: %q\n", err)
				}

				if stat.Code() != codes.FailedPrecondition {
					t.Fatalf("Got = %d, Want = %d\n", stat.Code(), codes.FailedPrecondition)
				}
			},
		},

		{
			name: "Create order with invalid market",
			f: func(t *testing.T) {
				req := client.CreateRequest{
					UserId:    userID,
					MarketId:  uuid.NewString(),
					OrderType: client.OrderType_ORDER_TYPE_T1,
					Price:     123,
					Quantity:  1,
				}
				_, err := orderService.Create(baseCtx, &req)

				if err == nil {
					t.Fatalf("Got nil error")
				}

				stat, ok := status.FromError(err)

				if !ok {
					t.Fatalf("Error on convert status from error: %q\n", err)
				}

				if stat.Code() != codes.FailedPrecondition {
					t.Fatalf("Got = %d, Want = %d\n", stat.Code(), codes.FailedPrecondition)
				}
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, tt.f)
	}
}

func TestStream(t *testing.T) {
	req := client.OrderUpdatesRequest{
		OrderId: orderId,
		UserId:  userID,
		DelayMs: (time.Millisecond * 300).Milliseconds(),
	}
	stream, err := orderService.OrderUpdates(baseCtx, &req)

	if err != nil {
		t.Fatalf("Got = %q\n", err)
	}

	var wantStatus = client.OrderStatus_ORDER_STATUS_CREATED

	for {
		resp, err := stream.Recv()

		if err != nil {
			if err == io.EOF {
				break
			}

			t.Fatalf("Got = %q\n", err)
		}

		switch wantStatus {
		case client.OrderStatus_ORDER_STATUS_CREATED:
			if resp.GetStatus() != wantStatus {
				t.Fatalf("Got = %d, Want = %d\n", resp.GetStatus(), wantStatus)
			}
			wantStatus = client.OrderStatus_ORDER_STATUS_PROCESSING

		case client.OrderStatus_ORDER_STATUS_PROCESSING:
			if resp.GetStatus() != wantStatus {
				t.Fatalf("Got = %d, Want = %d\n", resp.GetStatus(), wantStatus)
			}
			wantStatus = client.OrderStatus_ORDER_STATUS_PROCESSED

		case client.OrderStatus_ORDER_STATUS_PROCESSED:
			if resp.GetStatus() != wantStatus {
				t.Fatalf("Got = %d, Want = %d\n", resp.GetStatus(), wantStatus)
			}
			wantStatus = client.OrderStatus_ORDER_STATUS_CONFIRM

		default:
			if resp.GetStatus() != client.OrderStatus_ORDER_STATUS_CONFIRM && resp.GetStatus() != client.OrderStatus_ORDER_STATUS_REJECT {
				t.Fatalf("Got = %d, Want = %d || %d\n", resp.GetStatus(), client.OrderStatus_ORDER_STATUS_CONFIRM, client.OrderStatus_ORDER_STATUS_REJECT)

			} else {
				err := stream.CloseSend()

				if err != nil {
					t.Fatalf("Got = %q\n", err)
				}

				break
			}
		}
	}
}
