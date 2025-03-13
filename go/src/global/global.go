package global

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// MongoDB连接
	MongoDB *mongo.Client
)
