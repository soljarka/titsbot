package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotUser struct {
	user       *tgbotapi.User
	lastActive time.Time
	ignored    bool
}

type BotUsers struct {
	userList map[int]*BotUser
}

func (this *BotUsers) Ignore(User *tgbotapi.User) {
	foundUser, ok := this.userList[User.ID]

	if !ok {
		this.userList[User.ID] = &BotUser{
			user:       User,
			lastActive: time.Now(),
			ignored:    true,
		}
	} else {
		foundUser.ignored = true
	}
}

func (this *BotUsers) CheckIgnore(User *tgbotapi.User) bool {
	foundUser, ok := this.userList[User.ID]

	if !ok {
		return true
	}

	if foundUser.ignored == true {
		return true
	}
	return false
}

func (this *BotUsers) CheckInterval(User *tgbotapi.User) bool {
	foundUser, ok := this.userList[User.ID]

	if !ok {
		return true
	}

	duration := time.Since(foundUser.lastActive)
	if duration.Seconds() < 5 {
		return false
	}
	return true
}

func (this *BotUsers) AddUser(User *tgbotapi.User) {
	this.userList[User.ID] = &BotUser{
		user:       User,
		lastActive: time.Now(),
		ignored:    false,
	}
}
