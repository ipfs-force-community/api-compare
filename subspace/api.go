package main

import (
	"context"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/venus/venus-shared/api"
)

type SlotInfo struct {
	// Slot number
	SlotNumber uint64
	// Global slot challenge
	GlobalChallenge []byte
	// Acceptable solution range for block authoring
	SolutionRange uint64
	// Acceptable solution range for voting
	VotingSolutionRange uint64
}

type RewardSigningInfo struct {
	// Hash to be signed.
	Hash []byte
	// Public key of the plot identity that should create signature.
	PublicKey []byte
}

type SubspaceAPI interface {
	SubscribeSlotInfo(ctx context.Context) (<-chan *SlotInfo, error)
	SubscribeRewardSigning(ctx context.Context) (<-chan *RewardSigningInfo, error)
}

type SubspaceAPIStruct struct {
	Internal struct {
		// GetFarmerAppInfo       func(ctx context.Context, height abi.ChainEpoch, tsk types.TipSetKey) (*types.TipSet, error) `rpc_method:"subspace_getFarmerAppInfo"`
		SubscribeSlotInfo      func(ctx context.Context) (<-chan *SlotInfo, error)          `rpc_method:"subspace_subscribeSlotInfo"`
		SubscribeRewardSigning func(ctx context.Context) (<-chan *RewardSigningInfo, error) `rpc_method:"subspace_subscribeRewardSigning"`
	}
}

func (s *SubspaceAPIStruct) SubscribeSlotInfo(ctx context.Context) (<-chan *SlotInfo, error) {
	return s.Internal.SubscribeSlotInfo(ctx)
}

func (s *SubspaceAPIStruct) SubscribeRewardSigning(ctx context.Context) (<-chan *RewardSigningInfo, error) {
	return s.Internal.SubscribeRewardSigning(ctx)
}

func dialSubspaceRPC(ctx context.Context, addr string, token string, requestHeader http.Header) (SubspaceAPI, jsonrpc.ClientCloser, error) {
	if requestHeader == nil {
		requestHeader = http.Header{}
	}

	var res SubspaceAPIStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "", api.GetInternalStructs(&res), requestHeader)

	return &res, closer, err
}
