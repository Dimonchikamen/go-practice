package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := fixArgs(os.Args[1:])
	paths, isSuccess := validateArgs(args)
	if isSuccess {
		if value, ok := paths["-f"]; ok {
			for i := 0; i < len(value); i++ {
				var file *os.File
				var img image.Image
				var err error

				ext := filepath.Ext(value[i])
				if strings.Contains(ext, ".jpg") || strings.Contains(ext, ".jpeg") {
					file = openFileWithoutClose(value[i])
					img, err = jpeg.Decode(file)
					if err != nil {
						panic(err)
					}
				} else if ext == ".png" {
					file = openFileWithoutClose(value[i])
					img, err = png.Decode(file)
					if err != nil {
						panic(err)
					}
				} else {
					continue
				}

				defer file.Close()
				distImg := markImage(img)
				saveImage(file, ext, distImg)
			}
		}

		// if value, ok := paths["-d"]; ok {

		// }
	}
}

func saveImage(file *os.File, ext string, markedImage image.Image) {
	out, err := os.Create(strings.TrimSuffix(file.Name(), ext) + "_marked" + ext)
	if err != nil {
		panic(err)
	}

	if strings.Contains(ext, ".jpg") || strings.Contains(ext, ".jpeg") {
		jpeg.Encode(out, markedImage, nil)
	} else if ext == ".png" {
		png.Encode(out, markedImage)
	}

	defer out.Close()
}

func markImage(srcImage image.Image) *image.RGBA {
	srcBounds := srcImage.Bounds()
	// Создаем текстовый водяной знак
	watermark := image.NewRGBA(image.Rect(0, 0, srcBounds.Dx()/2, srcBounds.Dy()/20))
	watermarkBounds := watermark.Bounds()
	draw.Draw(watermark, watermarkBounds, &image.Uniform{color.NRGBA{0, 0, 0, 128}}, image.ZP, draw.Src)

	// Накладываем водяной знак на изображение
	offset := image.Pt(0, 0)
	distImg := image.NewRGBA(srcBounds)
	draw.Draw(distImg, srcBounds, srcImage, image.ZP, draw.Over)
	draw.Draw(distImg, watermarkBounds.Add(offset), watermark, image.ZP, draw.Over)
	return distImg
}

func openFileWithoutClose(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

func getPaths(args []string, startIndex int) ([]string, int) {
	subRes := []string{}
	startIndex++
	for startIndex < len(args) && args[startIndex][0] != '-' {
		subRes = append(subRes, args[startIndex])
		startIndex++
	}
	return subRes, startIndex - 1
}

func setPaths(res []string, mapper map[string][]string, key string) {
	if _, ok := mapper[key]; !ok {
		mapper[key] = res
	} else {
		mapper[key] = append(mapper[key], res...)
	}
}

func fixArgs(args []string) []string {
	result := []string{}
	subRes := ""
	for i := 0; i < len(args); i++ {
		if strings.Contains(args[i], "/") || strings.Contains(args[i], "\\") {
			subRes += args[i] + " "
		} else if len(subRes) > 0 {
			result = append(result, subRes)
			subRes = ""
		} else {
			result = append(result, args[i])
		}
	}
	if len(subRes) > 0 {
		result = append(result, subRes)
	}
	return result
}

func validateArgs(args []string) (map[string][]string, bool) {
	result := make(map[string][]string)
	println(strings.Join(args, " "))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d", "--directory":
			subRes, newIndex := getPaths(args, i)
			i = newIndex
			setPaths(subRes, result, "-d")
		case "-f", "--file":
			subRes, newIndex := getPaths(args, i)
			i = newIndex
			setPaths(subRes, result, "-f")
		case "-h", "--help":
			displayHelp()
			return result, false
		default:
			println()
			println("Команда" + "\"" + args[i] + "\" не поддерживается")
			displayHelp()
			return result, false
		}
	}
	return result, true
}

func displayHelp() {
	println()
	println("Данное консольное приложение добавляет водный знак на изображения форматов png, jpg и jpeg, остальные форматы изображения и файлы будут проигнорированы")
	println("Доступные команды:")
	println("-d, --directory - Добавить пудь до директории, которую будет просматривать утилита и искать в ней изображения для сохранения водного знака")
	println("Пример использования: -d C:/myDocuments, C:/users/user/desktop/myPhotos")
	println("-f, --file - Добавить пудь до изображения(-ий), в которое(-ые) будет добавлен водный знак")
	println("Пример использования: -f C:/myDocuments/my_photo.jpg, C:/users/user/desktop/myPhotos/anyPhoto.png")
	println("-h, --help - Показать доступные команды и их применение")
	println()
}
