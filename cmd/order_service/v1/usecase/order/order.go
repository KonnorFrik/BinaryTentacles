package order

import (
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
	return o
}
