package model

type Connection struct {
	SourceComponentId     *int64  `json:"sourceComponentId"`
	SourcePointName       *string `json:"sourcePointName"`
	TargetComponentId     *int64  `json:"targetComponentId"`
	RelativePointPosition *Point  `json:"relativePointPosition"`
}
