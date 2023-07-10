package model

type CommandStatus int

var (
	StatusCommandActive CommandStatus
	StatusCommandDel    CommandStatus = 1
)

type Commands []*Command

type Command struct {
	Id          *int64        `json:"id"`
	Type        *string       `json:"type"`
	Data        *string       `json:"data"`
	ComponentId *int64        `json:"componentId"`
	NextStepId  *int64        `json:"nextStepId"`
	Status      CommandStatus `json:"-"`
}

type AddCommandReq struct {
	CommandParams
}

type UpdCommandReq struct {
	CommandParams
}

type CommandsParam []*CommandParams

type CommandParams struct {
	Type *string `json:"type"`
	Data *string `json:"data"`
}

type SetNextStepCommandReq struct {
	NextStepId *int64 `json:"nextStepId"`
}
