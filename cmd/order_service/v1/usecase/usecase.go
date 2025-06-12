package usecase

import (
	"context"
	"errors"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase/order"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	db "github.com/KonnorFrik/BinaryTentacles/pkg/fake_db"
)

// TODO: create market client and call it

var fakeDB = db.New()

var (
	ErrDoesNotExist = errors.New("object does not exist")
)

func Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (
	*order.Order,
	error,
) {
	var order = new(order.Order)
	order.FromGrpcCreateRequest(req)
	// TODO: validate order object before create
	id := fakeDB.Create(ctx, order)
	order.Id = id
	return order, nil
}

func OrderStatus(
	ctx context.Context,
	req *pb.OrderStatusRequest,
) (
	*order.Order,
	error,
) {
	order, ok := db.As[*order.Order](ctx, fakeDB, req.GetOrderId())

	if !ok {
		return nil, ErrDoesNotExist
	}

	return order, nil
}

func Get(
	ctx context.Context,
	id uint64,
) (
	*order.Order,
	error,
) {
	order, ok := db.As[*order.Order](ctx, fakeDB, id)

	if !ok {
		return nil, ErrDoesNotExist
	}

	return order, nil
}
