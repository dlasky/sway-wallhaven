package main

// SwayOutputs struct representing sways json format for current outputs
type SwayOutputs []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rect struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"rect"`
	Focus              []int   `json:"focus"`
	Border             string  `json:"border"`
	CurrentBorderWidth int     `json:"current_border_width"`
	Layout             string  `json:"layout"`
	Orientation        string  `json:"orientation"`
	Percent            float64 `json:"percent"`
	WindowRect         struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window_rect"`
	DecoRect struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"deco_rect"`
	Geometry struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"geometry"`
	Window           interface{}   `json:"window"`
	Urgent           bool          `json:"urgent"`
	FloatingNodes    []interface{} `json:"floating_nodes"`
	Sticky           bool          `json:"sticky"`
	Type             string        `json:"type"`
	Active           bool          `json:"active"`
	Primary          bool          `json:"primary"`
	Make             string        `json:"make"`
	Model            string        `json:"model"`
	Serial           string        `json:"serial"`
	Scale            float64       `json:"scale"`
	Transform        string        `json:"transform"`
	CurrentWorkspace string        `json:"current_workspace"`
	Modes            []struct {
		Width   int `json:"width"`
		Height  int `json:"height"`
		Refresh int `json:"refresh"`
	} `json:"modes"`
	CurrentMode struct {
		Width   int `json:"width"`
		Height  int `json:"height"`
		Refresh int `json:"refresh"`
	} `json:"current_mode"`
	Focused bool `json:"focused"`
}
