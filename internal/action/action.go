package action

// Action は pager に対する操作を表す。
// 入力方法（キー・マウス等）とは切り離して定義する。
type ActionType int

const (
	ActNone ActionType = iota
	ActDown
	ActUp
	ActQuit
	ActPageDown
	ActPageUp
	ActTop
	ActBottom
	ActToggleRuler
	ActIncreaseHeader
	ActDecreaseHeader
	ActIncreaseFooter
	ActDecreaseFooter
	ActIncreaseHeight
	ActDecreaseHeight
	ActScrollLeft
	ActScrollRight
)

type Action struct {
	Type  ActionType
	Count int
}
