package touch

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	previousTouchIDs touchIDsMap
)

type touchIDsMap struct {
	mu sync.Mutex
	m  map[int]TouchInfo
}

type TouchInfo struct {
	X int
	Y int
}

func (tm *touchIDsMap) Set(key int, value TouchInfo) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.m[key] = value
}

func (tm *touchIDsMap) IsExist(key int) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	_, ok := tm.m[key]

	return ok
}

func (tm *touchIDsMap) All() map[int]TouchInfo {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	copy := make(map[int]TouchInfo)
	for k, v := range tm.m {
		copy[k] = v
	}

	return copy
}

func (tm *touchIDsMap) Reset() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.m = make(map[int]TouchInfo)
}

func IsTouchJustReleased() (bool, *TouchInfo) {
	result := false
	var touchInfo TouchInfo

	var currentTouchIDs []ebiten.TouchID
	currentTouchIDs = ebiten.AppendTouchIDs(currentTouchIDs)

	// タッチ開始を検出する
	for _, id := range currentTouchIDs {
		if exists := previousTouchIDs.IsExist(int(id)); !exists {
			// ...
		}
	}

	// タッチ終了を検出する
	for id, info := range previousTouchIDs.All() {
		found := false
		for _, currentID := range currentTouchIDs {
			if id == int(currentID) {
				found = true
				break
			}
		}
		if !found {
			result = true
			touchInfo = info
		}
	}

	// 現在のタッチIDを次回のために保存する
	previousTouchIDs.Reset()
	for _, id := range currentTouchIDs {
		x, y := ebiten.TouchPosition(id)
		previousTouchIDs.Set(int(id), TouchInfo{X: x, Y: y})
	}

	return result, &touchInfo
}
