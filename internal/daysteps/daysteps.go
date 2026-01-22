package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (steps int, duration time.Duration, err error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return steps, duration, errors.New("неверный формат данных")
	}

	// Количество шагов
	steps, err = strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, errors.New("неверное количество шагов")
	}

	// Продолжительность прогулки
	duration, err = time.ParseDuration(parts[1])
	if err != nil || duration <= 0 {
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	return
}

func DayActionInfo(data string, weight, height float64) (info string) {
	// Шаги и продолжительность прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("ошибка при разборе данных: %v", err)
		return
	}

	// Дистанция в метрах
	distanceMeters := float64(steps) * stepLength
	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	// Сожжённые калории
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	// Результат
	info = fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)

	return
}
