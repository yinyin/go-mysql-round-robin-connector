package mysqlroundrobinconnector

import (
	"log"
	"strconv"
	"sync"
)

type parsedAddress struct {
	locationName string
	orderedCount int
	shuffleCount int
}

var (
	addrCacheLock sync.RWMutex
	addrCache     map[string]*parsedAddress
)

func checkAddrCache(addrText string) (parsedAddr *parsedAddress) {
	addrCacheLock.RLock()
	defer addrCacheLock.RUnlock()
	/*
		if addrCache == nil {
			return
		}
	*/
	parsedAddr = addrCache[addrText]
	return
}

func updateAddrCache(addrText string, parsedAddr *parsedAddress) {
	addrCacheLock.Lock()
	defer addrCacheLock.Unlock()
	if addrCache == nil {
		addrCache = make(map[string]*parsedAddress)
	}
	addrCache[addrText] = parsedAddr
}

func parseAddress(addrText string) (parsedAddr *parsedAddress) {
	if parsedAddr = checkAddrCache(addrText); parsedAddr != nil {
		return
	}
	locNameIdx := 0
	orderedCntIdx := 0
	shuffleCntIdx := 0
	boundIdx := len(addrText)
	for idx, ch := range addrText {
		if ch != '/' {
			continue
		}
		if locNameIdx == 0 {
			locNameIdx = idx
			orderedCntIdx = boundIdx
		} else if orderedCntIdx == boundIdx {
			orderedCntIdx = idx
			shuffleCntIdx = boundIdx
		} else if shuffleCntIdx == boundIdx {
			shuffleCntIdx = idx
		}
	}
	if locNameIdx == 0 {
		locNameIdx = boundIdx
	}
	if locNameIdx == 0 {
		return nil
	}
	locName := addrText[:locNameIdx]
	orderedCnt := -1
	if orderedCntIdx > 0 {
		t := addrText[locNameIdx+1 : orderedCntIdx]
		if (t != "") && (t != "-") {
			v, err := strconv.ParseUint(t, 10, 16)
			if nil != err {
				log.Printf("ERROR: cannot parse ordered address count [%s] for location %s: %v", t, locName, err)
				return nil
			}
			orderedCnt = int(v)
		}
	}
	shuffleCnt := -1
	if shuffleCntIdx > 0 {
		t := addrText[orderedCntIdx+1 : shuffleCntIdx]
		if (t != "") && (t != "-") {
			v, err := strconv.ParseUint(t, 10, 16)
			if nil != err {
				log.Printf("ERROR: cannot parse shuffle address count [%s] for location %s: %v", t, locName, err)
				return nil
			}
			shuffleCnt = int(v)
		}
	}
	parsedAddr = &parsedAddress{
		locationName: locName,
		orderedCount: orderedCnt,
		shuffleCount: shuffleCnt,
	}
	updateAddrCache(addrText, parsedAddr)
	return
}
