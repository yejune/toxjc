package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yejune/toxjc/lib"
)

func TestConvertCSVtoJSON(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	jsonPath := filepath.Join(tmpDir, "test.json")

	if err := os.WriteFile(csvPath, []byte("a,b\n1,2"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := toxjc.Convert(csvPath, jsonPath); err != nil {
		t.Errorf("Convert failed: %v", err)
	}

	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}

func TestConvertJSONtoCSV(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "test.json")
	csvPath := filepath.Join(tmpDir, "test.csv")

	if err := os.WriteFile(jsonPath, []byte(`[{"a":"1","b":"2"}]`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := toxjc.Convert(jsonPath, csvPath); err != nil {
		t.Errorf("Convert failed: %v", err)
	}

	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}

func TestConvertCSVtoXLSX(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")

	if err := os.WriteFile(csvPath, []byte("a,b\n1,2"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := toxjc.Convert(csvPath, xlsxPath); err != nil {
		t.Errorf("Convert failed: %v", err)
	}

	if _, err := os.Stat(xlsxPath); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}

func TestConvertJSONtoXLSX(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "test.json")
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")

	if err := os.WriteFile(jsonPath, []byte(`[{"a":"1","b":"2"}]`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := toxjc.Convert(jsonPath, xlsxPath); err != nil {
		t.Errorf("Convert failed: %v", err)
	}

	if _, err := os.Stat(xlsxPath); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}
