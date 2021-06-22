package http

import "math/rand"

type ResponseHandler struct {
	random    bool // flag for random mode
	totalResp int  // total Responses available in imposters

	counter     int   // to keep count of served requests (wrapping after totalResp)
	currentInd  int   // index/key of current response in scheduleMap
	scheduleMap []int // prefix array of repeating request
}

// fillDefaults populates values based on imposter configuration
func (tr *ResponseHandler) fillDefaults(imposter *Imposter) {
	var scheduleMap = make([]int, len(imposter.Responses))
	burstAvailable := false

	for ind, resp := range imposter.Responses {
		if resp.Burst != 0 {
			burstAvailable = true
		} else {
			resp.Burst = 1
		}

		if ind != 0 {
			scheduleMap[ind] = scheduleMap[ind-1] + resp.Burst
		} else {
			scheduleMap[ind] = resp.Burst
		}
	}

	tr.scheduleMap = scheduleMap
	tr.counter = 1
	tr.currentInd = 0
	tr.totalResp = len(imposter.Responses)

	if !burstAvailable {
		tr.random = true
	}
}

// GetIndex is responsible for getting index for current request
func (dr *ResponseHandler) GetIndex() int {
	if dr.random {
		return dr.getRandomIndex()
	}
	return dr.getDynaimcIndex()
}

// getRandomIndex generates random indexes in random mode
func (dr *ResponseHandler) getRandomIndex() int {
	return rand.Intn(dr.totalResp)
}

// getDynamicIndex generates dynamic index based on the config provided
func (dr *ResponseHandler) getDynaimcIndex() int {
	var index = dr.currentInd

	dr.counter += 1 // incrementing counter for current request

	// checking if it has to move to next response or not
	if dr.scheduleMap[dr.currentInd] < dr.counter {
		dr.currentInd += 1
	}

	// Wrapping logic for counter and index
	if dr.currentInd > dr.totalResp-1 {
		dr.currentInd = 0
		dr.counter = 1
	}

	return index // returning current request
}
