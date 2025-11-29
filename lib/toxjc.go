// Package toxjc provides file format conversion utilities (to CSV/JSON/XLSX).
// Supports CSV, XLSX, XLS, JSON with automatic file type detection.
package toxjc

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// =============== Internal Helpers ===============

func detectEncoding(data []byte) string {
	if bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		return "utf-8-bom"
	}
	if bytes.HasPrefix(data, []byte{0xFF, 0xFE}) {
		return "utf-16-le"
	}
	if bytes.HasPrefix(data, []byte{0xFE, 0xFF}) {
		return "utf-16-be"
	}
	if utf8.Valid(data) {
		return "utf-8"
	}
	return "euc-kr"
}

func toUTF8Reader(r io.Reader, encoding string) io.Reader {
	switch encoding {
	case "utf-8", "utf-8-bom":
		return r
	case "utf-16-le":
		return transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder())
	case "utf-16-be":
		return transform.NewReader(r, unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder())
	case "euc-kr":
		return transform.NewReader(r, korean.EUCKR.NewDecoder())
	default:
		return r
	}
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

// =============== Detect ===============

// Detect: 파일 시그니처로 실제 타입 감지 (csv, xlsx, xls, json만 지원)
func Detect(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	header := make([]byte, 8)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return "", err
	}
	header = header[:n]

	// XLSX: ZIP (PK)
	if n >= 2 && header[0] == 0x50 && header[1] == 0x4B {
		return "xlsx", nil
	}

	// XLS: OLE2
	if n >= 8 && header[0] == 0xD0 && header[1] == 0xCF &&
		header[2] == 0x11 && header[3] == 0xE0 {
		return "xls", nil
	}

	// 텍스트 파일 체크 (바이너리면 not supported)
	file.Seek(0, 0)
	buf := make([]byte, 4096)
	n, _ = file.Read(buf)

	// 널 바이트 있으면 바이너리
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return "", fmt.Errorf("not supported")
		}
	}

	content := strings.TrimSpace(string(buf[:n]))
	if len(content) == 0 {
		return "", fmt.Errorf("not supported")
	}

	// JSON: [ 또는 { 로 시작
	if content[0] == '[' || content[0] == '{' {
		return "json", nil
	}

	// CSV: 콤마 또는 탭 구분자가 있는 텍스트
	if strings.Contains(content, ",") || strings.Contains(content, "\t") {
		return "csv", nil
	}

	return "", fmt.Errorf("not supported")
}

// =============== Read ===============

// ReadCSV: CSV 파일을 레코드로 읽기
func ReadCSV(path string) ([][]string, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("파일 열기 실패: %v", err)
	}
	defer inFile.Close()

	header := make([]byte, 4096)
	n, _ := inFile.Read(header)
	encoding := detectEncoding(header[:n])
	inFile.Seek(0, 0)

	var reader io.Reader = inFile
	if encoding == "utf-8-bom" {
		inFile.Seek(3, 0)
	} else if encoding == "utf-16-le" || encoding == "utf-16-be" {
		inFile.Seek(2, 0)
	}
	reader = toUTF8Reader(inFile, encoding)

	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1

	var records [][]string
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		for i, field := range record {
			record[i] = sanitize(field)
		}
		records = append(records, record)
	}
	return records, nil
}

// ReadXLSX: XLSX 파일을 레코드로 읽기
func ReadXLSX(path string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("XLSX 열기 실패: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("시트가 없음")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("행 읽기 실패: %v", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("빈 시트")
	}

	headerCols := len(rows[0])
	var records [][]string
	for _, row := range rows {
		for i, cell := range row {
			row[i] = sanitize(cell)
		}
		for len(row) < headerCols {
			row = append(row, "")
		}
		records = append(records, row)
	}
	return records, nil
}

// ReadXLS: XLS 파일을 레코드로 읽기
func ReadXLS(path string) ([][]string, error) {
	xlFile, err := xls.Open(path, "utf-8")
	if err != nil {
		return ReadXLSX(path) // fallback
	}

	sheet := xlFile.GetSheet(0)
	if sheet == nil {
		return nil, fmt.Errorf("시트가 없음")
	}

	headerCols := 0
	if firstRow := sheet.Row(0); firstRow != nil {
		headerCols = firstRow.LastCol() + 1
	}

	var records [][]string
	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			continue
		}
		var cells []string
		for j := 0; j <= row.LastCol(); j++ {
			cells = append(cells, sanitize(row.Col(j)))
		}
		for len(cells) < headerCols {
			cells = append(cells, "")
		}
		records = append(records, cells)
	}
	return records, nil
}

// ReadJSON: JSON 파일을 레코드로 읽기
func ReadJSON(path string) ([][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("파일 읽기 실패: %v", err)
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("JSON 파싱 실패: %v", err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("빈 JSON 배열")
	}

	var header []string
	for key := range items[0] {
		header = append(header, key)
	}

	records := [][]string{header}
	for _, item := range items {
		row := make([]string, len(header))
		for i, key := range header {
			if val, ok := item[key]; ok {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		records = append(records, row)
	}
	return records, nil
}

// Read: 파일 타입 자동 감지 후 읽기
func Read(path string) ([][]string, error) {
	fileType, err := Detect(path)
	if err != nil {
		return nil, fmt.Errorf("파일 타입 감지 실패: %v", err)
	}

	switch fileType {
	case "csv":
		return ReadCSV(path)
	case "xlsx":
		return ReadXLSX(path)
	case "xls":
		return ReadXLS(path)
	case "json":
		return ReadJSON(path)
	default:
		return nil, fmt.Errorf("지원하지 않는 파일 형식: %s", fileType)
	}
}

// =============== Write ===============

// WriteCSV: 레코드를 CSV로 저장
func WriteCSV(records [][]string, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("출력 파일 생성 실패: %v", err)
	}
	defer outFile.Close()

	w := csv.NewWriter(outFile)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			return fmt.Errorf("CSV 쓰기 실패: %v", err)
		}
	}
	return nil
}

// WriteJSON: 레코드를 JSON으로 저장
func WriteJSON(records [][]string, path string) error {
	if len(records) < 1 {
		return fmt.Errorf("데이터가 없음")
	}

	header := records[0]
	var data []map[string]string
	for _, row := range records[1:] {
		item := make(map[string]string)
		for i, col := range header {
			if i < len(row) {
				item[col] = row[i]
			} else {
				item[col] = ""
			}
		}
		data = append(data, item)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("출력 파일 생성 실패: %v", err)
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

// WriteXLSX: 레코드를 XLSX로 저장
func WriteXLSX(records [][]string, path string) error {
	if len(records) < 1 {
		return fmt.Errorf("데이터가 없음")
	}

	f := excelize.NewFile()
	sheet := "Sheet1"

	for rowIdx, record := range records {
		for colIdx, cell := range record {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			f.SetCellValue(sheet, cellName, cell)
		}
	}
	return f.SaveAs(path)
}

// =============== Convenience Converters ===============

// ToCSV: 파일을 CSV로 변환
func ToCSV(in, out string) error {
	records, err := Read(in)
	if err != nil {
		return err
	}
	if err := WriteCSV(records, out); err != nil {
		return err
	}
	fmt.Printf("  변환 완료: %s (%d rows)\n", out, len(records)-1)
	return nil
}

// ToJSON: 파일을 JSON으로 변환
func ToJSON(in, out string) error {
	records, err := Read(in)
	if err != nil {
		return err
	}
	if err := WriteJSON(records, out); err != nil {
		return err
	}
	fmt.Printf("  변환 완료: %s (%d rows)\n", out, len(records)-1)
	return nil
}

// ToXLSX: 파일을 XLSX로 변환
func ToXLSX(in, out string) error {
	records, err := Read(in)
	if err != nil {
		return err
	}
	if err := WriteXLSX(records, out); err != nil {
		return err
	}
	fmt.Printf("  변환 완료: %s (%d rows)\n", out, len(records)-1)
	return nil
}

// Convert: 출력 확장자로 형식 결정하여 변환
func Convert(in, out string) error {
	ext := strings.ToLower(filepath.Ext(out))
	switch ext {
	case ".csv":
		return ToCSV(in, out)
	case ".json":
		return ToJSON(in, out)
	case ".xlsx":
		return ToXLSX(in, out)
	default:
		return fmt.Errorf("지원하지 않는 출력 형식: %s (지원: csv, json, xlsx)", ext)
	}
}
