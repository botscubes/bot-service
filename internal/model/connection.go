package model

type SourceConnectionPoint struct {
	SourceComponentId *int64  `json:"sourceComponentId"`
	SourcePointName   *string `json:"sourcePointName"`
}

type ConnectionPoint struct {
	SourceConnectionPoint
	RelativePointPosition *Point `json:"relativePointPosition"`
}

type Connection struct {
	ConnectionPoint

	TargetComponentId *int64 `json:"targetComponentId"`
}
