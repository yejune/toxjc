package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yejune/toxjc"
)

func printUsage() {
	fmt.Println("toxjc - 파일 형식 변환 도구")
	fmt.Println()
	fmt.Println("사용법:")
	fmt.Println("  toxjc <input> <output>    파일 변환 (출력 확장자로 형식 결정)")
	fmt.Println("  toxjc detect <file>       파일 타입 감지")
	fmt.Println()
	fmt.Println("지원 형식:")
	fmt.Println("  입력: csv, xlsx, xls, json (자동 감지)")
	fmt.Println("  출력: csv, xlsx, json")
	fmt.Println()
	fmt.Println("예시:")
	fmt.Println("  toxjc data.xlsx output.csv")
	fmt.Println("  toxjc data.csv output.json")
	fmt.Println("  toxjc data.json output.xlsx")
	fmt.Println("  toxjc detect unknown.dat")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	// detect 명령
	if cmd == "detect" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "사용법: toxjc detect <file>")
			os.Exit(1)
		}
		filePath := os.Args[2]
		fileType, err := toxjc.Detect(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "감지 실패: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(fileType)
		return
	}

	// 변환 명령
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// 입력 파일 존재 확인
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "파일을 찾을 수 없습니다: %s\n", inputPath)
		os.Exit(1)
	}

	// 출력 형식 확인
	outExt := strings.ToLower(filepath.Ext(outputPath))
	validExt := map[string]bool{".csv": true, ".xlsx": true, ".json": true}
	if !validExt[outExt] {
		fmt.Fprintf(os.Stderr, "지원하지 않는 출력 형식: %s\n", outExt)
		fmt.Fprintln(os.Stderr, "지원: .csv, .xlsx, .json")
		os.Exit(1)
	}

	// 변환
	if err := toxjc.Convert(inputPath, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "변환 실패: %v\n", err)
		os.Exit(1)
	}
}
