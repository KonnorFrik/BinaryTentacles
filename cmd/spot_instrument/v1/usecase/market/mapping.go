package market

import (
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/spot_instrument/v1"
)

// ToProtobuf - convert Market object in pb.Market object.
func ToProtobuf(in *Market, out *pb.Market) {
	if in == nil || out == nil {
		return
	}

	in.mut.Lock()
	out.Id = in.Id
	in.mut.Unlock()
}

// ToProtobufMany - convert many Market objects in pb.Market objects.
func ToProtobufMany(in []*Market, out []*pb.Market) {
	if in == nil {
		return
	}

	maxInd := max(len(in), len(out))

	for i := range maxInd {
		out[i] = new(pb.Market)
		ToProtobuf(in[i], out[i])
	}
}
