package controllers

import "provider_gateway_mq/schema"

// ReplyChannels keep waiting channels for reply messages from rabbit
var ReplyChannels = make(map[string]chan schema.ReplyMessage)

// PublishChannels provide channels to publish rabbit messages
var PublishChannels = make(chan schema.CreateMessage, 10)
