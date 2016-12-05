package model

type EventJson struct {
	Timestamp           string `json:"@timestamp"`
	EventName           string `json:"eventName"`
	EventSourceAddress  string `json:"eventSourceAddress"`
	ProductVersion      string `json:"productVersion"`
	Counter             int `json:"counter"`
	GlobalTransactionId string `json:"globalTransactionId"`
	LocalTransactionId  string `json:"localTransactionId"`
	ParentTransactionId string `json:"parentTransactionId"`
	BrokerName          string `json:"brokerName"`
	HostName          string `json:"brokerName"`
	ExecutionGroupName  string `json:"executionGroupName"`
	MessageFlowName     string `json:"messageFlowName"`
	UniqueFlowName      string `json:"uniqueFlowName"`
	NodeDetail          string `json:"nodeDetail"`
	NodeType            string `json:"nodeType"`
	NodeLabel           string `json:"nodeLabel"`
	NodeTerminal        string `json:"nodeTerminal"`
	SimpleContents      map[string]string `json:"simpleContents"`
	ComplexContents     []string `json:"complexContents"`
	BitstreamEncoding   string `json:"bitstreamEncoding"`
	Bitstream           string `json:"bitstream"`
}
