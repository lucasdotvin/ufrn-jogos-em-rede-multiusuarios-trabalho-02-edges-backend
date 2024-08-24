package broadcast

type OnDisconnect func(channelKey string, subscriberKey string)

type Callbacker struct {
	OnDisconnect
}
