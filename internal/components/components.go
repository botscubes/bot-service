package components

// TODO: create wrapper
// component
// 	+ type
// 	+ content
// 	.. etc
//
// ComponentWrapEd
// 	- component
// 		+ type
// 		+ content
// 		 .. etc
//  - position
//
// 	ComponentWrapBot
//  - component
//

type Component struct {
	// *int64 may be worse than int64
	// BUT
	// ! *int64 used for optional fields !

	Id         *int64      `json:"id"`
	Data       *Data       `json:"data"`
	Keyboard   *Keyboard   `json:"keyboard"`
	Commands   *[]*Command `json:"commands"`
	NextStepId *int64      `json:"nextStepId"`
	IsMain     bool        `json:"isMain"`
	Position   *Point      `json:"position"`
}

type Command struct {
	Id          *int64  `json:"id,omitempty"`
	Type        *string `json:"type"`
	Data        *string `json:"data"`
	ComponentId *int64  `json:"componentId"`
	NextStepId  *int64  `json:"nextStepId"`
}

type Point struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
}

type Data struct {
	Type    *string     `json:"type"`
	Content *[]*Content `json:"content"`
}

type Content struct {
	Text *string `json:"text,omitempty"`
}

type Keyboard struct {
	Buttons [][]*int64 `json:"buttons"`
}
