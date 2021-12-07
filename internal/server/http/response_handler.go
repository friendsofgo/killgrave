package http

import (
	"math/rand"
)

// ResponseMode represents random/burst mode for the response
type ResponseMode int

// ResponseHandler handles incoming request for random/burst response
type ResponseHandler struct {
	Mode      ResponseMode // flag for random mode
	totalResp int          // total Responses available in imposters

	counter     int   // to keep count of served requests (wrapping after totalResp)
	currentInd  int   // index/key of current response in scheduleMap
	scheduleMap []int // prefix array of repeating request
}

const (
	// RandomMode will generate random responses
	RandomMode ResponseMode = iota
	// BurstMode will generate repeatable responses
	BurstMode
)

// fillDefaults populates values based on imposter configuration
func (rh *ResponseHandler) fillDefaults(imposter *Imposter) {
	// Updating totalResponse length
	rh.totalResp = len(imposter.Responses)

	// Updating Response Mode
	switch imposter.Request.ResponseMode {
	case "BURST":
		rh.Mode = BurstMode
	default:
		rh.Mode = RandomMode
	}

	// Populating state for BURST mode
	if rh.Mode == BurstMode {
		var scheduleMap = make([]int, len(imposter.Responses))
		for ind, resp := range imposter.Responses {
			value := resp.Burst
			if value <= 0 {
				value = 1
			}
			if ind != 0 {
				scheduleMap[ind] = scheduleMap[ind-1] + value
			} else {
				scheduleMap[ind] = value
			}
		}

		rh.scheduleMap = scheduleMap
		rh.counter = 1
		rh.currentInd = 0
	}
}

// GetIndex is responsible for getting index for current request
func (rh *ResponseHandler) GetIndex() int {
	if rh.Mode == RandomMode {
		return rh.getRandomIndex()
	}
	ind := rh.getBurstIndex()

	return ind
}

// getRandomIndex generates random indexes in random mode
func (rh *ResponseHandler) getRandomIndex() int {
	return rand.Intn(rh.totalResp)
}

// getBurstIndex generates repeated index based on the config provided
func (rh *ResponseHandler) getBurstIndex() int {
	var index = rh.currentInd

	rh.counter++ // incrementing counter for current request

	// checking if it has to move to next response or not
	if rh.scheduleMap[rh.currentInd] < rh.counter {
		rh.currentInd++
	}

	// Wrapping logic for counter and index
	if rh.currentInd > rh.totalResp-1 {
		rh.currentInd = 0
		rh.counter = 1
	}

	return index // returning current request
}
