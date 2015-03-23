// postioned object managment in 2d space
package idpos1m

import (
	"sync"

	"github.com/kasworld/idpos/idposi"
	// "github.com/kasworld/log"
)

type Manager struct {
	xlen, ylen int
	id2obj     map[int64]idposi.IDPosI
	id2posi    map[int64]int
	pos2objs   map[int]idposi.IDPosIList
	mutex      sync.Mutex
}

func (t *Manager) pos2Int(p [2]int) int {
	return t.xlen*p[1] + p[0]
}
func (t *Manager) posXY2Int(x, y int) int {
	return t.xlen*y + x
}

func New(x, y int) idposi.IDPosManI {
	rtn := Manager{
		xlen:     x,
		ylen:     y,
		id2obj:   make(map[int64]idposi.IDPosI),
		id2posi:  make(map[int64]int),
		pos2objs: make(map[int]idposi.IDPosIList),
		// pos2objs: make(map[int]idposi.IDPosIList, x*y),
	}
	return &rtn
}

func (fo *Manager) Count() int {
	return len(fo.id2obj)
}
func (fo *Manager) All() map[int64]idposi.IDPosI {
	return fo.id2obj
}

func (fo *Manager) GetByID(id int64) idposi.IDPosI {
	return fo.id2obj[id]
}

func (fo *Manager) addPos2Objs(o idposi.IDPosI, posi int) {
	for i, v := range fo.pos2objs[posi] {
		if v == nil {
			fo.pos2objs[posi][i] = o
			return
		}
	}
	fo.pos2objs[posi] = append(fo.pos2objs[posi], o)
}

func (fo *Manager) delPos2Objs(o idposi.IDPosI, posi int) bool {
	for i, v := range fo.pos2objs[posi] {
		if v == nil {
			continue
		}
		if v.GetID() == o.GetID() {
			fo.pos2objs[posi][i] = nil
			return true
		}
	}
	// critical error
	return false
}

func (fo *Manager) Set(o idposi.IDPosI) bool {
	if fo.id2obj[o.GetID()] != nil {
		if !fo.Del(o) {
			return false
		}
	}
	fo.Add(o)
	return true
}

func (fo *Manager) Add(o idposi.IDPosI) bool {
	id := o.GetID()
	pos := o.GetPos()
	posi := fo.pos2Int(pos)
	if fo.id2obj[id] != nil {
		// log.Error("obj exist %v %v", id, pos)
		return false
	}

	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	fo.id2obj[id] = o
	fo.id2posi[id] = posi
	fo.addPos2Objs(o, posi)
	return true
}

func (fo *Manager) Del(o idposi.IDPosI) bool {
	id := o.GetID()
	if fo.id2obj[id] == nil {
		// log.Error("obj not exist %v", id)
		return false
	}
	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	posi := fo.id2posi[id]
	if fo.delPos2Objs(o, posi) {
		delete(fo.id2obj, id)
		delete(fo.id2posi, id)
		return true
	}
	return false
}

func (fo *Manager) UpdateToPos(o idposi.IDPosI, newpos [2]int) bool {
	fo.mutex.Lock()
	defer fo.mutex.Unlock()
	oldpos := o.GetPos()
	if fo.delPos2Objs(o, fo.pos2Int(oldpos)) {
		newposi := fo.pos2Int(newpos)
		fo.id2posi[o.GetID()] = newposi
		fo.addPos2Objs(o, newposi)
		return true
	}
	return false
}

func (fo *Manager) iterAt(posi int, fn idposi.DoFn) bool {
	for _, v := range fo.pos2objs[posi] {
		if v == nil {
			continue
		}
		if fn(v) {
			return true
		}
	}
	return false
}
func (fo *Manager) IterAtXY(x, y int, fn idposi.DoFn) bool {
	posi := fo.posXY2Int(x, y)
	return fo.iterAt(posi, fn)
}

func (fo *Manager) IterAt(pos [2]int, fn idposi.DoFn) bool {
	posi := fo.pos2Int(pos)
	return fo.iterAt(posi, fn)
}

func (fo *Manager) PosXYObjs(x, y int) idposi.IDPosIList {
	return fo.pos2objs[fo.posXY2Int(x, y)]
}