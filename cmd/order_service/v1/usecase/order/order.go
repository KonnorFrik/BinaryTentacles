package order

import (
	"context"
	"math/rand"
	"sync"

	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
)

type Order struct {
	mut sync.Mutex

	Id       uint64
	MarketId uint64
	Type     pb.OrderType
	Price    float64
	Quantity uint64

	Status pb.OrderStatus
}

func (o *Order) GetStatus() pb.OrderStatus {
	o.mut.Lock()
	defer o.mut.Unlock()
	return o.Status
}

func (o *Order) FromGrpcCreateRequest(
	req *pb.CreateRequest,
) *Order {
	o.MarketId = req.GetMarketId()
	o.Type = req.GetOrderType()
	o.Price = req.GetPrice()
	o.Quantity = req.GetQuantity()
	return o
}

func (o *Order) ToGrpcCreateResponse(
	resp *pb.CreateResponse,
) *Order {
	resp.OrderId = o.Id
	resp.OrderStatus = o.Status
	return o
}

// UpdateStatus - Start goroutine for fake update order status.
func (o *Order) UpdateStatus(
	ctx context.Context,
) <-chan pb.OrderStatus {
	var result = make(chan pb.OrderStatus)

	go func() {
		var loop bool = true
		for loop {
			o.mut.Lock()
			var orderStatus pb.OrderStatus = o.Status

			if o.Status == pb.OrderStatus_ORDER_STATUS_CONFIRM || o.Status == pb.OrderStatus_ORDER_STATUS_REJECT {
				break
			}

			switch o.Status {
			case pb.OrderStatus_ORDER_STATUS_CREATED:
				o.Status = pb.OrderStatus_ORDER_STATUS_PROCESSING

			case pb.OrderStatus_ORDER_STATUS_PROCESSING:
				o.Status = pb.OrderStatus_ORDER_STATUS_PROCESSED

			case pb.OrderStatus_ORDER_STATUS_PROCESSED:
				if rand.Intn(2) == 0 {
					o.Status = pb.OrderStatus_ORDER_STATUS_CONFIRM

				} else {
					o.Status = pb.OrderStatus_ORDER_STATUS_REJECT
				}

				loop = false
			}

			o.mut.Unlock()

			select {
			case <-ctx.Done():
				return
			case result <- orderStatus:
			}
		}

		close(result)
	}()

	return result
}
