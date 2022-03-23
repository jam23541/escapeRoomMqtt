package escapeRoomMqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jam23541/escapeRoomDataTypes"
)

// ChanRawMsgFromMqtt is the channel that contains msgs that is not in the mqttMsg format
var ChanRawMsgFromMqtt = make(chan mqtt.Message, 1000)

// ChanCleanMsgFromMqtt is the channel that contains clean msg to this code
var ChanCleanMsgFromMqtt = make(chan escapeRoomDataTypes.MqttMsg, 1000)

// ChanMsgToPub is the channel that contains all msg to be published
var ChanMsgToPub = make(chan escapeRoomDataTypes.MqttMsg, 1000)

func StringChannelWriteIn(newMsg string, myChannel *chan string) int {
	channelCapacity := cap(*myChannel)
	currentChannelLength := len(*myChannel)
	if currentChannelLength >= channelCapacity {
		return -1
	} else {
		*myChannel <- newMsg
		return currentChannelLength + 1
	}
}
func StringChannelReadOut(myChannel *chan string) (string, int) {
	//channelCapacity := cap(*my_channel)
	currentChannelLength := len(*myChannel)
	if currentChannelLength <= 0 {
		return "", -1
	} else {
		outputMsg := <-*myChannel
		return outputMsg, currentChannelLength - 1
	}

}

func IntChannelWriteIn(newMsg int, myChannel *chan int) int {
	channelCapacity := cap(*myChannel)
	currentChannelLength := len(*myChannel)
	if currentChannelLength >= channelCapacity {
		return -1
	} else {
		*myChannel <- newMsg
		return currentChannelLength + 1
	}
}
func IntChannelReadOut(myChannel *chan int) (int, int) {
	//channelCapacity := cap(*my_channel)
	currentChannelLength := len(*myChannel)
	if currentChannelLength <= 0 {
		return -1, -1
	} else {
		outputMsg := <-*myChannel
		return outputMsg, currentChannelLength - 1
	}

}

// MqttMsgChannelWriteIn WriteIn returns -1 if the channel is full, otherwise return the length of channel after msg added
func MqttMsgChannelWriteIn(newMsg escapeRoomDataTypes.MqttMsg, myChannel *chan escapeRoomDataTypes.MqttMsg) int {
	channelCapacity := cap(*myChannel)
	currentChannelLength := len(*myChannel)
	if currentChannelLength >= channelCapacity {
		return -1
	} else {
		*myChannel <- newMsg
		return currentChannelLength + 1
	}
}

// MqttMsgChannelReadOut ReadOut return -1 if the channel is already empty, and set the output string to be "EMPTY";
//otherwise return the msg along with the length of channel after taking out the msg
func MqttMsgChannelReadOut(myChannel *chan escapeRoomDataTypes.MqttMsg) (escapeRoomDataTypes.MqttMsg, int) {
	//channelCapacity := cap(*my_channel)
	currentChannelLength := len(*myChannel)
	if currentChannelLength <= 0 {

		return escapeRoomDataTypes.DefaultMqttMsg, -1
	} else {
		outputMsg := <-*myChannel
		return outputMsg, currentChannelLength - 1
	}

}
