package controllers

import "provider_gateway_mq/schemas"

// ReplyChannels keep waiting channels for reply messages from rabbit
var ReplyChannels = make(map[string]chan schemas.MessageReply)

// PublishChannels channel to publish rabbit messages
var PublishChannels = make(chan schemas.MessageCreate, 10)
