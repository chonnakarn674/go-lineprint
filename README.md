# go-lineprinter

A Go library for **line matrix printer reports**, supporting Thai text alignment and legacy printer encoding.

---

## Features

* Fixed-position formatting (column layout)
* Thai alignment correction
* Windows-874 encoding support
* 3-level Thai rendering (upper/base/bottom)
* Buffer → Redis → Printer workflow

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

### 2. Render Thai text for line matrix printers 

```go
var buf bytes.Buffer

lineprint.Print(&buf, "กำลังทดสอบ")
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

### Print

```
Print(io.Writer, string)
```

---

## License

MIT License - see LICENSE file
