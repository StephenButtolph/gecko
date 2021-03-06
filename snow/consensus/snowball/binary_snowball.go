// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowball

import (
	"fmt"
)

// binarySnowball is the implementation of a binary snowball instance
type binarySnowball struct {
	// preference is the choice with the largest number of successful polls.
	// Ties are broken by switching choice lazily
	preference int

	// numSuccessfulPolls tracks the total number of successful network polls of
	// the 0 and 1 choices
	numSuccessfulPolls [2]int

	// snowflake wraps the binary snowflake logic
	snowflake binarySnowflake
}

// Initialize implements the BinarySnowball interface
func (sb *binarySnowball) Initialize(beta, choice int) {
	sb.preference = choice
	sb.snowflake.Initialize(beta, choice)
}

// Preference implements the BinarySnowball interface
func (sb *binarySnowball) Preference() int {
	// It is possible, with low probability, that the snowflake preference is
	// not equal to the snowball preference when snowflake finalizes. However,
	// this case is handled for completion. Therefore, if snowflake is
	// finalized, then our finalized snowflake choice should be preferred.
	if sb.Finalized() {
		return sb.snowflake.Preference()
	}
	return sb.preference
}

// RecordSuccessfulPoll implements the BinarySnowball interface
func (sb *binarySnowball) RecordSuccessfulPoll(choice int) {
	sb.numSuccessfulPolls[choice]++
	if sb.numSuccessfulPolls[choice] > sb.numSuccessfulPolls[1-choice] {
		sb.preference = choice
	}
	sb.snowflake.RecordSuccessfulPoll(choice)
}

// RecordUnsuccessfulPoll implements the BinarySnowball interface
func (sb *binarySnowball) RecordUnsuccessfulPoll() { sb.snowflake.RecordUnsuccessfulPoll() }

// Finalized implements the BinarySnowball interface
func (sb *binarySnowball) Finalized() bool { return sb.snowflake.Finalized() }

func (sb *binarySnowball) String() string {
	return fmt.Sprintf(
		"SB(Preference = %d, NumSuccessfulPolls[0] = %d, NumSuccessfulPolls[1] = %d, SF = %s)",
		sb.preference,
		sb.numSuccessfulPolls[0],
		sb.numSuccessfulPolls[1],
		&sb.snowflake)
}
