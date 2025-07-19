package storage

import "github.com/SHSanderland/EffMobTest/pkg/model"

type Storage interface {
	CreateSubscription(sub *model.Subscription) error
}
