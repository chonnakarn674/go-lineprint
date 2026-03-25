# go-lineprinter

A Go library for generating **line matrix printer reports**, with support for Thai text alignment and legacy printer encoding.

---

## Features

- Fixed-position formatting (column-based layout)
- Thai text alignment support
- Windows-874 (TIS-620) encoding
- 3-level Thai rendering (upper / base / bottom)
- Buffer → Redis → Printer workflow support

---

## Installation

```bash
go get github.com/chonnakarn674/go-lineprint
```

---

## Requirements

- Go 1.20+
- golang.org/x/text (auto-installed via Go modules)

---

## Usage

### 1. Format report layout

```go
var buf bytes.Buffer

lineprint.Format(&buf, []int{1, 25, 50}, "ชื่อ", "นามสกุล", "สถานะ")
lineprint.Format(&buf, []int{1, 25}, "สมชาย", "นามสมมุติ")

lineprint.FormatLine(&buf)
lineprint.FormatNewPage(&buf)
```

---

### 2. Convert Thai text into 3-level line printer format

```go
var buf bytes.Buffer

data := lineprint.ToLP3("กำลังทดสอบ")

buf.Write(data)
```

---

## Processing Pipeline

```text
Format → Buffer → Redis → Retrieve → Print
```

---

## Redis Example

```go
data := buf.Bytes()
redisClient.Set(ctx, "job", data, 0)

data, _ = redisClient.Get(ctx, "job").Bytes()
printer.Write(data)
```

---

## API

### Format

```
Format(io.Writer, []int, ...string)
```

### ToLP3

```
ToLP3(string) []byte
```

---

## License

MIT License - see LICENSE file
