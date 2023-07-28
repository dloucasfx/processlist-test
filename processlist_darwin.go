//go:build darwin
// +build darwin

package main

import "github.com/sirupsen/logrus"

type osCache struct {
}

func initOSCache() *osCache {
	return &osCache{}
}

func ProcessList(cache *osCache, logger logrus.FieldLogger) ([]*TopProcess, error) {
	return nil, nil
}
