package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Количество шагов
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("неверное количество шагов")
	}

	// Продолжительность прогулки
	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	return steps, duration, nil
}
func DayActionInfo(data string, weight, height float64) string {
	// Шаги и продолжительность прогулки
	steps, _, err := parsePackage(data)
	if err != nil {
		log.Printf("ошибка при разборе данных: %v", err)
		return ""
	}

	// Проверка количества шагов
	if steps <= 0 {
		return ""
	}

	// Дистанция в метрах
	distanceMeters := float64(steps) * stepLength
	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	// Сожжённые калории
	calories := WalkingSpentCalories(weight, height, distanceKm)

	// Результат
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)

	return result
}

func WalkingSpentCalories(weight, height, distanceKm float64) any {
	panic("unimplemented")
}
