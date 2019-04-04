package main

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	//据StackOverflow, cost=10在2018年的Mac笔记本上耗时75ms，基本
	//不影响用户体验
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
