package panels


import "elevators/core"

type OneButton struct {

}

func ToPanel(core.Floor, core.Floor) OneButton {
	return OneButton{}
}
