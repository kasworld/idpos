// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// postioned object managment in 2d space
package idpos2s

import (
	"sync"

	"github.com/kasworld/idgen"
	"github.com/kasworld/idpos/idposi"
)

type Manager struct {
	id2obj   map[idgen.IDInt]idposi.IDPosI
	id2pos   map[idgen.IDInt][2]int
	pos2objs [][]idposi.IDPosIList
	mutex    sync.RWMutex
}

func New(x, y int) idposi.IDPosManI {
	rtn := Manager{
		id2obj:   make(map[idgen.IDInt]idposi.IDPosI),
		id2pos:   make(map[idgen.IDInt][2]int),
		pos2objs: make([][]idposi.IDPosIList, x),
	}
	for i, _ := range rtn.pos2objs {
		rtn.pos2objs[i] = make([]idposi.IDPosIList, y)
	}
	return &rtn
}

func (fo *Manager) Count() int {
	return len(fo.id2obj)
}
func (fo *Manager) All() idposi.IDPosIList {
	rtn := make(idposi.IDPosIList, 0, len(fo.id2obj))
	fo.mutex.RLock()
	defer fo.mutex.RUnlock()
	for _, v := range fo.id2obj {
		rtn = append(rtn, v)
	}
	return rtn
}
func (fo *Manager) GetByID(id idgen.IDInt) idposi.IDPosI {
	return fo.id2obj[id]
}
func (fo *Manager) PosXYObjs(x, y int) idposi.IDPosIList {
	return fo.pos2objs[x][y]
}

func (fo *Manager) addPos2Objs(o idposi.IDPosI, x, y int) {
	for i, v := range fo.pos2objs[x][y] {
		if v == nil {
			fo.pos2objs[x][y][i] = o
			return
		}
	}
	fo.pos2objs[x][y] = append(fo.pos2objs[x][y], o)
}

func (fo *Manager) delPos2Objs(o idposi.IDPosI, x, y int) bool {
	for i, v := range fo.pos2objs[x][y] {
		if v == nil {
			continue
		}
		if v.GetID() == o.GetID() {
			fo.pos2objs[x][y][i] = nil
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
	if fo.id2obj[id] != nil {
		// log.Error("obj exist %v %v", id, pos)
		return false
	}

	fo.mutex.Lock()
	defer fo.mutex.Unlock()

	fo.id2obj[id] = o
	fo.id2pos[id] = pos
	fo.addPos2Objs(o, pos[0], pos[1])
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

	pos := fo.id2pos[id]
	if fo.delPos2Objs(o, pos[0], pos[1]) {
		delete(fo.id2obj, id)
		delete(fo.id2pos, id)
		return true
	}
	return false
}

func (fo *Manager) UpdateToPos(o idposi.IDPosI, newpos [2]int) bool {
	fo.mutex.Lock()
	defer fo.mutex.Unlock()
	oldpos := o.GetPos()
	if fo.delPos2Objs(o, oldpos[0], oldpos[1]) {
		fo.id2pos[o.GetID()] = newpos
		fo.addPos2Objs(o, newpos[0], newpos[1])
		return true
	}
	return false
}

func (fo *Manager) IterAtXY(x, y int, fn idposi.DoFn) bool {
	for _, v := range fo.pos2objs[x][y] {
		if v == nil {
			continue
		}
		if fn(v) {
			return true
		}
	}
	return false
}

func (fo *Manager) IterAt(pos [2]int, fn idposi.DoFn) bool {
	return fo.IterAtXY(pos[0], pos[1], fn)
}
