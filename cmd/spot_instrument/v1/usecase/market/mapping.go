package market

import (
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
)

func ToProtobuf(in *Market, out *pb.Market) {
	if in == nil || out == nil {
		return
	}

	in.mut.Lock()
	out.Id = in.Id
	in.mut.Unlock()
}

func ToProtobufMany(in []*Market, out []*pb.Market) {
	if in == nil {
		return
	}

	maxInd := max(len(in), len(out))

	for i := range maxInd {
		ToProtobuf(in[i], out[i])
	}
}
