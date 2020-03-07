package model

type Datapoint struct {
	Timestamp  int64   `json:"timestamp"`
	UserId     string  `json:"userId"`
	InstanceId string  `json:"instanceId"`
	Value      float64 `json:"Value"`
	// ValueFloat64 float64 `json:"-"`
	// just for copy
	// `json:""`
}

// ParseValue like 3.2468807808000002E7
// func (data *Datapoint) ParseValue() (err error) {
// 	number, err := strconv.ParseFloat(data.Value, 64)
// 	if err != nil {
// 		return err
// 	}
// 	data.ValueFloat64 = number
// 	return nil
// }
