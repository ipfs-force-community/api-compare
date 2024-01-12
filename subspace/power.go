package main

import (
	"github.com/docker/go-units"
)

const (
	max uint64 = 18446744073709551615

	// 1 in 6 slots (on average, not counting collisions) will have a block.
	// Must match ratio between block and slot duration in constants above.
	slotProbabilityOne uint64 = 1
	slotProbabilityTwo uint64 = 6

	maxPiecesInSector uint64 = 1000
	numChunks         uint64 = 32768
	numSBuckets       uint64 = 65536

	pieceSize uint64 = 1048672
)

// MAX * slot_probability / (pieces_in_sector * chunks / s_buckets) / sectors
func sectorsToSolutionRange(sectors uint64) uint64 {
	solutionRange := max / slotProbabilityTwo * slotProbabilityOne / (maxPiecesInSector * numChunks / numSBuckets)
	// Account for slot probability
	// Now take sector size and probability of hitting occupied s-bucket in sector into account

	// Take number of sectors into account
	return solutionRange / sectors
}

// MAX * slot_probability / (pieces_in_sector * chunks / s_buckets) / solution_range
func solutionRangeToSectors(solutionRange uint64) uint64 {
	sectors := max / slotProbabilityTwo * slotProbabilityOne / (maxPiecesInSector * numChunks / numSBuckets)
	// Account for slot probability
	// Now take sector size and probability of hitting occupied s-bucket in sector into account

	// Take solution range into account
	return sectors / solutionRange
}

func networkPower(solutionRange uint64) uint64 {
	sectors := solutionRangeToSectors(solutionRange)

	return sectors * maxPiecesInSector * pieceSize
}

func pretty(x uint64) string {
	return units.BytesSize(float64(x))
}
