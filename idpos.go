// postioned object managment in 2d space
package idpos

import (
	"sync"

	"github.com/kasworld/log"
)

type IDPosI interface {
	GetID() int64
	GetPos() [2]int
}

type Manager struct {
	id2obj   map[int64]IDPosI
	id2pos   map[int64][2]int
	pos2objs map[[2]int][]IDPosI
	mutex    sync.Mutex
}

func New() *Manager {
	rtn := Manager{
		id2obj:   make(map[int64]IDPosI),
		id2pos:   make(map[int64][2]int),
		pos2objs: make(map[[2]int][]IDPosI),
	}
	return &rtn
}

func (fo *Manager) Count() int {
	return len(fo.id2obj)
}
func (fo *Manager) All() map[int64]IDPosI {
	return fo.id2obj
}

func (fo *Manager) addPos2Objs(o IDPosI, pos [2]int) {
	for i, v := range fo.pos2objs[pos] {
		if v == nil {
			fo.pos2objs[pos][i] = o
			return
		}
	}
	fo.pos2objs[pos] = append(fo.pos2objs[pos], o)
}

func (fo *Manager) delPos2Objs(o IDPosI, pos [2]int) bool {
	for i, v := range fo.pos2objs[pos] {
		if v == nil {
			continue
		}
		if v.GetID() == o.GetID() {
			fo.pos2objs[pos][i] = nil
			return true
		}
	}
	// critical error
	return false
}

func (fo *Manager) Set(o IDPosI) bool {
	if fo.id2obj[o.GetID()] != nil {
		if !fo.Del(o) {
			return false
		}
	}
	fo.Add(o)
	return true
}

func (fo *Manager) Add(o IDPosI) bool {
	id := o.GetID()
	pos := o.GetPos()
	if fo.id2obj[id] != nil {
		// log.Error("obj exist %v %v", id, pos)
		return false
	}

	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	fo.id2obj[id] = o
	fo.id2pos[id] = pos
	fo.addPos2Objs(o, pos)
	return true
}

func (fo *Manager) Del(o IDPosI) bool {
	id := o.GetID()
	if fo.id2obj[id] == nil {
		// log.Error("obj not exist %v", id)
		return false
	}
	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	pos := fo.id2pos[id]
	if fo.delPos2Objs(o, pos) {
		delete(fo.id2obj, id)
		delete(fo.id2pos, id)
		return true
	}
	return false
}

func (fo *Manager) CheckPos(o IDPosI) bool {
	return fo.id2pos[o.GetID()] == o.GetPos()
}

func (fo *Manager) CheckAll() bool {
	rtn := true
	for _, v := range fo.All() {
		if !fo.CheckPos(v) {
			rtn = false
			log.Error("id2pos check fail %v %v", v.GetID(), v.GetPos())
		}
	}
	for p, _ := range fo.pos2objs {
		fo.IterAt(p, func(o IDPosI) bool {
			if !fo.CheckPos(o) {
				log.Error("pos2obj check fail %v %v", o.GetID(), o.GetPos())
				rtn = false
			}
			return false
		})
	}
	return rtn
}

func (fo *Manager) UpdateToPos(o IDPosI, newpos [2]int) bool {
	fo.mutex.Lock()
	defer fo.mutex.Unlock()
	oldpos := o.GetPos()
	if fo.delPos2Objs(o, oldpos) {
		fo.id2pos[o.GetID()] = newpos
		fo.addPos2Objs(o, newpos)
		return true
	}
	return false
}

type DoFn func(fo IDPosI) bool

func (fo *Manager) IterAtXY(x, y int, fn DoFn) bool {
	return fo.IterAt([2]int{x, y}, fn)
}

func (fo *Manager) IterAt(pos [2]int, fn DoFn) bool {
	for _, v := range fo.pos2objs[pos] {
		if v == nil {
			continue
		}
		if fn(v) {
			return true
		}
	}
	return false
}