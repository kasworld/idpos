package idposi

type IDPosI interface {
	GetID() int64
	GetPos() [2]int
}

type IDPosIList []IDPosI

type DoFn func(fo IDPosI) bool

type IDPosManI interface {
	Count() int
	All() map[int64]IDPosI
	GetByID(id int64) IDPosI
	IterAtXY(x, y int, fn DoFn) bool
	IterAt(pos [2]int, fn DoFn) bool
	PosXYObjs(x, y int) IDPosIList
	Set(o IDPosI) bool
	Add(o IDPosI) bool
	Del(o IDPosI) bool
	UpdateToPos(o IDPosI, newpos [2]int) bool
}
