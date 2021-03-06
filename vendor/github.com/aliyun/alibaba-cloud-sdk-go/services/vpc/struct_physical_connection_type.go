package vpc

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// PhysicalConnectionType is a nested struct in vpc response
type PhysicalConnectionType struct {
	PhysicalConnectionId          string `json:"PhysicalConnectionId" xml:"PhysicalConnectionId"`
	AccessPointId                 string `json:"AccessPointId" xml:"AccessPointId"`
	Type                          string `json:"Type" xml:"Type"`
	Status                        string `json:"Status" xml:"Status"`
	BusinessStatus                string `json:"BusinessStatus" xml:"BusinessStatus"`
	CreationTime                  string `json:"CreationTime" xml:"CreationTime"`
	EnabledTime                   string `json:"EnabledTime" xml:"EnabledTime"`
	LineOperator                  string `json:"LineOperator" xml:"LineOperator"`
	Spec                          string `json:"Spec" xml:"Spec"`
	PeerLocation                  string `json:"PeerLocation" xml:"PeerLocation"`
	PortType                      string `json:"PortType" xml:"PortType"`
	RedundantPhysicalConnectionId string `json:"RedundantPhysicalConnectionId" xml:"RedundantPhysicalConnectionId"`
	Name                          string `json:"Name" xml:"Name"`
	Description                   string `json:"Description" xml:"Description"`
	AdLocation                    string `json:"AdLocation" xml:"AdLocation"`
	PortNumber                    string `json:"PortNumber" xml:"PortNumber"`
	CircuitCode                   string `json:"CircuitCode" xml:"CircuitCode"`
	Bandwidth                     int64  `json:"Bandwidth" xml:"Bandwidth"`
	LoaStatus                     string `json:"LoaStatus" xml:"LoaStatus"`
	HasReservationData            string `json:"HasReservationData" xml:"HasReservationData"`
	ReservationInternetChargeType string `json:"ReservationInternetChargeType" xml:"ReservationInternetChargeType"`
	ReservationActiveTime         string `json:"ReservationActiveTime" xml:"ReservationActiveTime"`
	ReservationOrderType          string `json:"ReservationOrderType" xml:"ReservationOrderType"`
	EndTime                       string `json:"EndTime" xml:"EndTime"`
	ChargeType                    string `json:"ChargeType" xml:"ChargeType"`
}
