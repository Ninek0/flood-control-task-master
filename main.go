package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const N = 10
const K = 5

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

// Структура для хранения данных флуд-контроля
type floodControlData struct {
	mu          sync.Mutex
	callHistory map[int64][]time.Time
}

// Реализация метода Check для интерфейса
func (f *floodControlData) Check(ctx context.Context, userID int64) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	history := f.callHistory[userID]

	cutoffTime := time.Now().Add(-N * time.Second)
	var newHistory []time.Time
	for _, t := range history {
		if t.After(cutoffTime) {
			newHistory = append(newHistory, t)
		}
	}
	newHistory = append(newHistory, time.Now())
	f.callHistory[userID] = newHistory

	if len(newHistory) > K {
		return false, nil
	}

	return true, nil
}

func NewFloodControl() FloodControl {
	return &floodControlData{
		callHistory: make(map[int64][]time.Time),
	}
}

func main() {

	floodControl := NewFloodControl()

	ctx := context.TODO()
	userID := int64(1337)
	result, err := floodControl.Check(ctx, userID)
	if err != nil {
		fmt.Println("Ошибка при проверке флуд-контроля:", err.Error())
		return
	}
	if !result {
		fmt.Printf("Флуд-контроль для пользователя %d не пройден\n", userID)
	} else {
		fmt.Printf("Флуд-контроль для пользователя %d пройден\n", userID)
	}

	userID2 := int64(321)

	for i := 0; i < 10; i++ {
		result2, err2 := floodControl.Check(ctx, userID2)
		if err2 != nil {
			fmt.Println("Ошибка при проверке флуд-контроля:", err2.Error())
			return
		}
		if !result2 {
			fmt.Printf("Флуд-контроль для пользователя %d не пройден\n", userID2)
			return
		} else {
			fmt.Printf("Флуд-контроль для пользователя %d пройден\n", userID2)
		}
	}
}
