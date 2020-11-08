package main

import (
	"errors"
	"log"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	DistributedLockPrefix  = "LOCK:" // prefix of lock-key
	DistributedLockTimeout = 5       // timeout for lock, unit:second
)

func Lock(lockType string) (string, error) {
	if conn == nil {
		log.Fatal("Please init redis connection first.")
	}

	strLockName := strings.Join([]string{DistributedLockPrefix, lockType}, "")
	strUUID := uuid.NewV4().String()
	iTimeBegin := time.Now().Unix()
	for {
		replySetNx, err := conn.Do("SETNX", strLockName, strUUID)
		if err != nil {
			return "", err
		}
		if replySetNx.(int64) == 1 {
			_, _ = conn.Do("EXPIRE", strLockName, DistributedLockTimeout)
			return strUUID, nil
		}

		replyTTL, err := conn.Do("TTL", strLockName)
		if err != nil {
			return "", err
		}
		if replyTTL.(int64) == -1 {
			_, _ = conn.Do("EXPIRE", strLockName, DistributedLockTimeout)
		}

		if time.Now().Unix()-iTimeBegin > DistributedLockTimeout {
			return "", errors.New("Operation timeout, please try again.")
		}
		time.Sleep(time.Microsecond * 1)
	}
}

func Unlock(lockType string, strUUID string) error {
	strLockName := strings.Join([]string{DistributedLockPrefix, lockType}, "")
	replyGet, err := conn.Do("GET", strLockName)
	if err != nil {
		return err
	}

	if string(replyGet.([]uint8)) != strUUID {
		return errors.New("Currently someone else has the lock, you cannot unlock it.")
	}

	_, _ = conn.Do("WATCH", strLockName)
	_, _ = conn.Do("MULTI")
	_, _ = conn.Do("DEL", strLockName)
	replyExec, err := conn.Do("EXEC")
	if err != nil {
		return err
	}
	if replyExec.([]interface{})[0].(int64) == 0 {
		return errors.New("The lock is currently unlocked by someone else.")
	}
	return nil
}
