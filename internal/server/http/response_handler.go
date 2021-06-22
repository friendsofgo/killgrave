package http

import "math/rand"

type ResponseMode int

type ResponseHandler struct {
	Mode      ResponseMode // flag for random mode
	totalResp int          // total Responses available in imposters

	counter     int   // to keep count of served requests (wrapping after totalResp)
	currentInd  int   // index/key of current response in scheduleMap
	scheduleMap []int // prefix array of repeating request
}

const (
	// JSONImposter allows to know when we're dealing with a JSON imposter
	RandomMode ResponseMode = iota
	// YAMLImposter allows to know when we're dealing with a YAML imposter
	BurstMode
)

// fillDefaults populates values based on imposter configuration
func (tr *ResponseHandler) fillDefaults(imposter *Imposter) {
	// Updating totalResponse length
	tr.totalResp = len(imposter.Responses)

	// Updating Response Mode
	switch imposter.Request.ResponseMode {
	case "BURST":
		tr.Mode = BurstMode
	default:
		tr.Mode = RandomMode
	}

	// Populating state for BURST mode
	if tr.Mode == BurstMode {
		var scheduleMap = make([]int, len(imposter.Responses))
		for ind, resp := range imposter.Responses {
			if ind != 0 {
				scheduleMap[ind] = scheduleMap[ind-1] + resp.Burst
			} else {
				scheduleMap[ind] = resp.Burst
			}
		}

		tr.scheduleMap = scheduleMap
		tr.counter = 1
		tr.currentInd = 0
	}
}

// GetIndex is responsible for getting index for current request
func (dr *ResponseHandler) GetIndex() int {
	if dr.Mode == RandomMode {
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
