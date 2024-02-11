package model

type ConnectionPoint struct {
	SourceComponentId     *int64  `json:"sourceComponentId"`
	SourcePointName       *string `json:"sourcePointName"`
	RelativePointPosition *Point  `json:"relativePointPosition"`
}

type Connection struct {
	ConnectionPoint

	TargetComponentId *int64 `json:"targetComponentId"`
}

type DelConnectionReq struct {
	SourceComponentId *int64  `json:"sourceComponentId"`
	SourcePointName   *string `json:"sourcePointName"`
	TargetComponentId *int64  `json:"targetComponentId"`
}
