package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов(не используется по тз lenStep).
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку на части по запятой
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных тренировки")
	}

	// Количество шагов
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("неверное количество шагов")
	}

	// Вид активности
	activityType := strings.TrimSpace(parts[1])
	if activityType == "" {
		return 0, "", 0, errors.New("вид активности не указан")
	}

	// Продолжительность
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// Длина одного шага по росту
	stepLength := height * stepLengthCoefficient

	// Общая дистанция в метрах
	totalDistanceMeters := float64(steps) * stepLength

	// Переводим в километры
	distanceKm := totalDistanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	// Продолжительность > 0
	if duration <= 0 {
		return 0
	}

	// Дистанция в километрах
	distanceKm := distance(steps, height)

	// Продолжительность в часы
	durationHours := duration.Hours()

	// Средняя скорость = дистанция / время (км/ч)
	speed := distanceKm / durationHours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Данные тренировки
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Дистанция и средняя скорость
	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64

	// Вид тренировки и считаем калории
	switch activityType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}

	// Строка результата
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activityType,
		duration.Hours(),
		distanceKm,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// Рассчитываем среднюю(км/ч)
	speed := meanSpeed(steps, height, duration)

	// Продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Количество калорий
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// Средняя(км/ч)
	speed := meanSpeed(steps, height, duration)

	// Продолжительность в минутах
	durationMinutes := duration.Minutes()

	// Рассчет калорий
	calories := (weight * speed * durationMinutes) / minInH

	// Коэффициент для ходьбы
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
