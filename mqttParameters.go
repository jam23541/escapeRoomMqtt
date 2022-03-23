package escapeRoomMqtt

type MqttParameters struct {
	MyTopic            string
	DebugTopic         string
	BackendTopic       string
	RemotePublishTopic string
	ClientName         string
	Puber              string
	HostAddr           string
	Port               int
}
