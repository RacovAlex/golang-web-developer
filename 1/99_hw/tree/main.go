package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printTree(out, path, "", printFiles)
}

func printTree(out io.Writer, path string, prefix string, printFiles bool) error {
	// открываем директорию
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// читаем содержимое папки
	files, err := file.Readdir(0)
	if err != nil {
		return err
	}

	// фильтруем файлы/директории
	var items []os.FileInfo
	for _, file := range files {
		if !printFiles && !file.IsDir() {
			continue
		}
		items = append(items, file)
	}

	// сортируем файлы/директории
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	// проходим по всем элементам и выводим их
	for i, file := range items {
		var line string
		if i == len(items)-1 {
			line = "└───"
		} else {
			line = "├───"
		}

		// если элемент - директория, выводим и то проходим по ней рекурсивно
		if file.IsDir() {
			fmt.Fprintf(out, "%s%s%s\n", prefix, line, file.Name())
			newPrefix := prefix + "│\t"
			if i == len(items)-1 {
				newPrefix = prefix + "\t"
			}
			err = printTree(out, filepath.Join(path, file.Name()), newPrefix, printFiles)
			if err != nil {
				return err
			}
		} else {
			// обрабатываем файл
			fileSize := file.Size()
			sizeStr := fmt.Sprintf(" (%db)", fileSize)
			if fileSize == 0 {
				sizeStr = " (empty)"
			}
			fmt.Fprintf(out, "%s%s%s%s\n", prefix, line, file.Name(), sizeStr)
		}
	}
	return nil
}
