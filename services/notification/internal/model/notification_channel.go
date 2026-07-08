package model

type NotificationChannel string

const (

	ChannelEmail NotificationChannel = "EMAIL"

	ChannelSMS NotificationChannel = "SMS"

	ChannelPush NotificationChannel = "PUSH"

)

func (c NotificationChannel) IsSupported() bool {

	switch c {

	case ChannelEmail,
		ChannelSMS,
		ChannelPush:

		return true

	default:

		return false

	}

}