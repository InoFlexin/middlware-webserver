package base

import (
	"errors"
	"net"
)

var sessionMap map[string]net.Conn = map[string]net.Conn{}

func GetSessions() ([]string, []net.Conn) {
	length := len(sessionMap)
	keys := make([]string, 0, length)
	values := make([]net.Conn, 0, length)

	for k, v := range sessionMap {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

func RemoveSession(id string) {
	val := sessionMap[id]

	if val != nil {
		delete(sessionMap, id)
	}
}

func AddSession(id string, conn net.Conn) (string, error) {
	val := sessionMap[id]

	if val != nil {
		return id, errors.New(id + " is already added")
	}
	sessionMap[id] = conn

	return id, nil
}

func GetSession(id string) (net.Conn, error) {
	val, exists := sessionMap[id]

	if !exists {
		return nil, errors.New(id + " is not found")
	}

	return val, nil
}
