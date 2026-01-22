package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (steps int, activityType string, duration time.Duration, err error) {
	// Разделяем строку на части по запятой
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("не полные входные данные")
	}

	// Количество шагов
	steps, err = strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("неверное количество шагов")
	}

	// Вид активности
	activityType = parts[1]
	if activityType == "" {
		return 0, "", 0, errors.New("вид активности не указан")
	}

	// Продолжительность
	duration, err = time.ParseDuration(parts[2])
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	return
}

func distance(steps int, height float64) (distanceKm float64) {
	// Длина одного шага по росту
	stepLength := height * stepLengthCoefficient

	// Общая дистанция в метрах
	totalDistanceMeters := float64(steps) * stepLength

	// Переводим в километры
	distanceKm = totalDistanceMeters / mInKm

	return
}

func meanSpeed(steps int, height float64, duration time.Duration) (speed float64) {
	// Продолжительность > 0
	if duration <= 0 {
		return
	}

	// Дистанция в километрах
	distanceKm := distance(steps, height)

	// Средняя скорость = дистанция / время (км/ч)
	speed = distanceKm / duration.Hours()

	return
}

func TrainingInfo(data string, weight, height float64) (info string, err error) {
	// Данные тренировки
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
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
		return "", fmt.Errorf("неизвестный тип тренировки\n")
	}

	// Строка результата
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		distanceKm,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (call float64, err error) {
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

	// Количество калорий
	call = (weight * speed * duration.Minutes()) / minInH

	return
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (calories float64, err error) {
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

	// Рассчет калорий
	calories = ((weight * speed * duration.Minutes()) / minInH) * walkingCaloriesCoefficient

	return
}
