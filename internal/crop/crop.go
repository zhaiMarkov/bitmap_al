package crop

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"bitmap/config"
	"bitmap/internal/core"
)

// Структура для хранения параметров кропа
type CropValues struct {
	OffSetX int
	OffSetY int
	Width   int
	Height  int
}

// Основная функция для обработки crop
func HandleCrop(b *core.BitMap) {
	height, width := b.GetDimensions()

	// Если флаги для кропа не указаны — выходим
	if len(config.CropFlag) == 0 {
		return
	}

	// Берем первый флаг для обработки
	flag := config.CropFlag[0]
	config.CropFlag = config.CropFlag[1:] // Убираем обработанный флаг

	// Проверяем количество -
	dashCount := strings.Count(flag, "-")

	if dashCount != 1 && dashCount != 3 {
		fmt.Println("Формат заданных данных неправильна")
		os.Exit(1)
	}

	// Парсим флаг --crop
	cleanedSlice, err := parseFlags(flag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cropValues CropValues

	// Определяем логику обработки флага по количеству значений
	switch len(cleanedSlice) {
	case 2:
		// Если указаны только OffSetX и OffSetY, ширина и высота будут максимальными до конца изображения
		cropValues.OffSetX, cropValues.OffSetY = handleTwoValues(cleanedSlice)
		cropValues.Width = int(width) - cropValues.OffSetX
		cropValues.Height = int(height) - cropValues.OffSetY
	case 4:
		// Если указаны все 4 значения
		cropValues.OffSetX, cropValues.OffSetY, cropValues.Width, cropValues.Height = handleFourValues(cleanedSlice)
	default:
		fmt.Println("Некорректное количество аргументов для флага crop")
		os.Exit(1)
	}

	// Валидация значений
	if !validate(cropValues, int(width), int(height)) {
		fmt.Println("Ошибка: неправильные значения crop")
		os.Exit(1)
	}

	// Нарезаем изображение по заданным значениям
	cropImage(b, cropValues)

	// Обновляем текущие размеры изображения после кропа
	height, width = b.GetDimensions()
}

// Парсинг флагов и очистка значений
func parseFlags(flag string) ([]string, error) {
	preparedData := strings.Split(flag, "-")
	var cleanedSlice []string
	for _, str := range preparedData {
		if str != "" {
			_, err := strconv.Atoi(str)
			if err != nil {
				return nil, fmt.Errorf("некорректное значение: %s", str)
			}
			cleanedSlice = append(cleanedSlice, str)
		}
	}
	return cleanedSlice, nil
}

// Обработка двух значений (OffSetX и OffSetY)
func handleTwoValues(cleanedSlice []string) (int, int) {
	values := convertToInt(cleanedSlice)
	OffSetX, OffSetY := values[0], values[1]
	return OffSetX, OffSetY
}

// Обработка четырех значений (OffSetX, OffSetY, Width и Height)
func handleFourValues(cleanedSlice []string) (int, int, int, int) {
	values := convertToInt(cleanedSlice)
	return values[0], values[1], values[2], values[3]
}

// Валидация параметров кропа
func validate(cropValues CropValues, width, height int) bool {
	return cropValues.OffSetX >= 0 && cropValues.OffSetY >= 0 &&
		cropValues.Width > 0 && cropValues.Height > 0 &&
		cropValues.OffSetX+cropValues.Width <= width &&
		cropValues.OffSetY+cropValues.Height <= height
}

// Конвертируем строки в int
func convertToInt(s []string) []int {
	var values []int
	for _, data := range s {
		val, _ := strconv.Atoi(data)
		values = append(values, val)
	}
	return values
}

// Нарезаем изображение по заданным координатам
func cropImage(b *core.BitMap, cropValues CropValues) {
	OffSetX, OffSetY, Width, Height := cropValues.OffSetX, cropValues.OffSetY, cropValues.Width, cropValues.Height

	pixels := b.GetPixels()
	croppedPixels := make([][]*core.Pixel, Height)

	for i := range croppedPixels {
		croppedPixels[i] = make([]*core.Pixel, Width)
	}

	// Проходим по пикселям, начиная с верхнего левого угла и идем вниз
	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			// Извлекаем пиксели начиная с верхней границы (OffSetY) и слева направо (OffSetX)
			croppedPixels[i][j] = pixels[len(pixels)-OffSetY-i-1][OffSetX+j]
		}
	}

	for i := 0; i < Height/2; i++ {
		for j := 0; j < Width; j++ {
			croppedPixels[i][j], croppedPixels[Height-i-1][j] = croppedPixels[Height-i-1][j], croppedPixels[i][j]
		}
	}

	// Устанавливаем новые размеры изображения
	b.SetDimensions(int32(Height), int32(Width))

	// Получаем текущий размер изображения
	originalImageSize := b.GetImageSize()

	// Вычисляем новый размер изображения
	newImageSize := uint32(Height * Width * 3)

	// Вычитаем разницу из общего размера файла
	newFileSize := b.GetFileSize() - (originalImageSize - newImageSize)

	// Устанавливаем новые размеры
	b.SetImageSize(newImageSize)
	b.SetFileSize(newFileSize)

	// Сохраняем нарезанные пиксели
	b.SetPixels(croppedPixels)
}
