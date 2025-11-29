# toxjc

**toxjc** - 파일 형식 변환 도구 (to XLS/JSON/CSV)

CSV, JSON, XLSX, XLS 파일 형식을 자동 감지하여 상호 변환합니다.

> 이름 "toxjc"의 유래: **to** **X**(LS) **J**(SON) **C**(SV)

## 설치

### 방법 1: Homebrew (권장)

```bash
brew install yejune/tap/toxjc
```

### 방법 2: go install 사용

```bash
go install github.com/yejune/toxjc@latest
toxjc install
```

## 사용법

### 파일 변환

```bash
toxjc <입력> <출력>
```

출력 파일의 확장자로 형식이 결정됩니다.

### 예시

```bash
# Excel → CSV
toxjc data.xlsx output.csv

# CSV → JSON
toxjc data.csv output.json

# JSON → Excel
toxjc data.json output.xlsx

# 구버전 Excel
toxjc legacy.xls output.csv
```

### 파일 타입 감지

```bash
toxjc detect <파일>
```

```bash
$ toxjc detect unknown.dat
csv
```

### 버전 확인

```bash
toxjc version
```

## 지원 형식

| 형식 | 입력 | 출력 |
|------|------|------|
| CSV  | ✅   | ✅   |
| XLSX | ✅   | ✅   |
| XLS  | ✅   | ❌   |
| JSON | ✅   | ✅   |

## 기능

- 파일 타입 자동 감지
- 인코딩 자동 감지 (UTF-8, UTF-16, EUC-KR)
- BOM (Byte Order Mark) 처리
- 간단한 CLI 인터페이스

## 라이선스

MIT License
