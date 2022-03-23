package escapeRoomMqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jam23541/escapeRoomDataTypes"
)

type MqttClient struct {
	Client      mqtt.Client
	MyParameter MqttParameters
}

func myOPTS(myClientID string, myUsername string, myPassword string, myBroker string, myPort int) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", myBroker, myPort))
	opts.SetClientID(myClientID)
	opts.SetUsername(myUsername)
	opts.SetPassword(myPassword)
	opts.SetDefaultPublishHandler(onMessageReceived)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	return opts
}

//when connection established
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

// when connection lost
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var onMessageReceived mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	channelCapacity := cap(ChanRawMsgFromMqtt)
	currentChannelLength := len(ChanRawMsgFromMqtt)
	if currentChannelLength >= channelCapacity {
		fmt.Println("no room for raw mqtt msg")
	} else {
		// two parallel channels, raw and clean
		newMqttMsg, deJsonErr := escapeRoomDataTypes.PayloadToMqttMsg(msg.Payload())
		if deJsonErr == 1 { // clean msg
			ChanCleanMsgFromMqtt <- newMqttMsg
		} else { // raw msg
			ChanRawMsgFromMqtt <- msg
		}
	}
}

// create a new mqtt client
func NewMqttClient(myParameter MqttParameters) *MqttClient {
	return &MqttClient{
		Client: mqtt.NewClient(myOPTS(myParameter.ClientName,
			myParameter.ClientName,
			myParameter.ClientName,
			myParameter.HostAddr,
			myParameter.Port)),
		MyParameter: myParameter,
	}
}
func (myClient *MqttClient) Pub(myTopic string, inputString string) bool {
	text := fmt.Sprintf(inputString)
	token := myClient.Client.Publish(myTopic, 0, false, text)
	token.Wait()
	return true
}

func (myClient *MqttClient) DebugLog(inputString string) bool {
	text := fmt.Sprintf(inputString)
	token := myClient.Client.Publish(myClient.MyParameter.DebugTopic, 0, false, text)
	token.Wait()
	return true
}

func (myClient *MqttClient) Sub(myTopic string) bool {
	topic := fmt.Sprintf("%s/#", myTopic)
	token := myClient.Client.Subscribe(topic, 0, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
	return true
}

func (myClient *MqttClient) Init() {
	if token := myClient.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (myClient *MqttClient) MqttPublishRoutine() {
	for {
		// publishing msg
		if len(ChanMsgToPub) > 0 {
			newMsgToPub, err := MqttMsgChannelReadOut(&ChanMsgToPub)
			if err >= 0 {
				newMsgInJson, MarshalErr := json.Marshal(newMsgToPub)
				if MarshalErr != nil {
					fmt.Println("Encoding Json failed")
				} else {
					myClient.Pub(newMsgToPub.PUBLISHTOPIC, string(newMsgInJson))
				}
			}
		}
	}
}
