package messages

type RequestMessage struct {
	startdate string `json:"startdate"`
	enddate   string `json:"enddate"`
}
