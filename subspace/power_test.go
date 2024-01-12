package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkPower(t *testing.T) {
	assert.Equal(t, solutionRangeToSectors(sectorsToSolutionRange(1)), uint64(1))
	assert.Equal(t, solutionRangeToSectors(sectorsToSolutionRange(2)), uint64(2))
	assert.Equal(t, solutionRangeToSectors(sectorsToSolutionRange(3)), uint64(3))

	solutionRanges := []uint64{1, 2, 3, 683396345, 654294934}
	expects := []uint64{
		6148914691021189888,  // 5.333EiB
		12297829381841034752, // 10.67EiB
		2049638229990839296,  // 1.778EiB
		9435512311104000,     // 8.38PiB
		9855180358784000,     // 8.753PiB
	}
	for idx, solutionRange := range solutionRanges {
		power := networkPower(solutionRange)
		fmt.Println(solutionRange, power, pretty(power))
		assert.Equal(t, expects[idx], power)
	}
}
