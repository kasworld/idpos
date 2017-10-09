package idposi

import (
	"github.com/kasworld/idgen"
)

type IDPosI interface {
	GetID() idgen.IDInt
	GetPos() [2]int
}

type IDPosIList []IDPosI

type DoFn func(fo IDPosI) bool

type IDPosManI interface {
	Count() int
	All() IDPosIList
	GetByID(id idgen.IDInt) IDPosI
	IterAtXY(x, y int, fn DoFn) bool
	IterAt(pos [2]int, fn DoFn) bool
	PosXYObjs(x, y int) IDPosIList
	Set(o IDPosI) bool
	Add(o IDPosI) bool
	Del(o IDPosI) bool
	UpdateToPos(o IDPosI, newpos [2]int) bool
}
