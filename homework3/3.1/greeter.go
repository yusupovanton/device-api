package greeter

import (
	"errors"
	"fmt"
	"strings"
)

func Greet(name string, hour int) (string, error) {
	//Сделал так, чтобы значения времени меньше нуля и больше 23 возвращали обрабатываемую ошибку, а не панику.
	greeting := "Good night"

	trimmedName := strings.Trim(name, " ")
	if hour >= 6 && hour < 12 {
		greeting = "Good morning"
	} else if hour >= 12 && hour < 18 {
		// Тут я исправил <= 18 что по сути создавало ошибочный кейс, где 18 по присваивалось и вечеру и дню
		greeting = "Hello"
	} else if hour >= 18 && hour < 22 {
		greeting = "Good evening"
	} else if hour < 0 || hour >= 24 {
		greeting = "Hello"
		return fmt.Sprintf("%s %s!", greeting, strings.Title(trimmedName)), errors.New("Wrong hour is feeded!")
	}

	return fmt.Sprintf("%s %s!", greeting, strings.Title(trimmedName)), nil
}
