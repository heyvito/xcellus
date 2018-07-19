package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tealeg/xlsx"
)

func main() {}

type worksheet struct {
	Title   string          `json:"title"`
	Headers []string        `json:"headers"`
	Rows    [][]interface{} `json:"rows"`
}

var instances = map[string]*xlsx.File{}

const invalidHandleErr = "Invalid handle. Please file a bug.\x00"

//export go_xcellus_process
func go_xcellus_process(data *C.char) (failed bool, err *C.char, len int, buf *C.char) {
	synthesizedData := C.GoString(data)
	var wb []*worksheet
	goErr := json.Unmarshal([]byte(synthesizedData), &wb)

	if goErr != nil {
		return true, errorToCString(goErr), 0, nil
	}

	file := xlsx.NewFile()

	for _, sheet := range wb {
		sheetPtr, goErr := file.AddSheet(sheet.Title)
		if goErr != nil {
			return true, errorToCString(goErr), 0, nil
		}

		addHeaders(sheetPtr, sheet)
		addRows(sheetPtr, sheet)
	}

	tmpBuffer := &bytes.Buffer{}
	writer := bufio.NewWriter(tmpBuffer)
	goErr = file.Write(writer)

	if goErr != nil {
		return true, errorToCString(goErr), 0, nil
	}

	len = tmpBuffer.Len()
	p := C.malloc(C.size_t(len))

	// copy the data into the buffer, by converting it to a Go array
	cBuf := (*[1 << 30]byte)(p)
	copy(cBuf[:], tmpBuffer.Bytes())

	failed = false
	err = nil
	buf = (*C.char)(p) // Ownership of P will be changed to Ruby's after
	// being received from the counterpart C code.
	// DO NOT free IT.
	return
}

//export go_xcellus_load
func go_xcellus_load(rawPath *C.char) (failed bool, err, handle *C.char) {
	path := C.GoString(rawPath)
	xl, goErr := xlsx.OpenFile(path)
	if goErr != nil {
		return true, errorToCString(goErr), nil
	}
	id, goErr := randomHex(20)
	if goErr != nil {
		return true, errorToCString(goErr), nil
	}
	instances[id] = xl
	return false, nil, C.CString(fmt.Sprintf("%s\x00", id))
}

//export go_xcellus_append
func go_xcellus_append(rawHandle, rawData *C.char) (failed bool, err *C.char) {
	data := C.GoString(rawData)
	handle := C.GoString(rawHandle)

	file, ok := instances[handle]
	if !ok {
		return true, C.CString(invalidHandleErr)
	}

	var wb []*worksheet
	goErr := json.Unmarshal([]byte(data), &wb)

	if goErr != nil {
		return true, errorToCString(goErr)
	}

	for _, b := range wb {
		sheetPtr := file.Sheet[b.Title]
		if sheetPtr == nil {
			sheetPtr, goErr := file.AddSheet(b.Title)
			if goErr != nil {
				return true, errorToCString(goErr)
			}
			addHeaders(sheetPtr, b)
			addRows(sheetPtr, b)
		} else {
			addRows(sheetPtr, b)
		}
	}

	return false, nil
}

//export go_xcellus_save
func go_xcellus_save(rawHandle, rawPath *C.char) (failed bool, err *C.char) {
	path := C.GoString(rawPath)
	handle := C.GoString(rawHandle)

	ioutil.WriteFile("/tmp/dat1", []byte(fmt.Sprintf("Received handle %s, has handles: %#v\n", handle, instances)), 0644)

	file, ok := instances[handle]
	if !ok {
		return true, C.CString(invalidHandleErr)
	}

	goErr := file.Save(path)
	if goErr != nil {
		return true, errorToCString(goErr)
	}

	return false, nil
}

//export go_xcellus_end
func go_xcellus_end(rawHandle *C.char) (failed bool, err *C.char) {
	handle := C.GoString(rawHandle)

	_, ok := instances[handle]
	if !ok {
		return true, C.CString(invalidHandleErr)
	}

	delete(instances, handle)

	return false, nil
}

//export go_xcellus_find_in_column
func go_xcellus_find_in_column(rawHandle, rawSheetName, rawValue *C.char, index int) (failed bool, err *C.char, idxResult int) {
	handle := C.GoString(rawHandle)
	sheetName := C.GoString(rawSheetName)
	value := C.GoString(rawValue)

	file, ok := instances[handle]
	if !ok {
		return false, C.CString(invalidHandleErr), -1
	}

	sheetPtr := file.Sheet[sheetName]
	if sheetPtr == nil { // Sheet does not exist. Move along.
		return false, nil, -1
	}

	return false, nil, findInColumn(sheetPtr, index, value)
}

//export go_xcellus_replace_row
func go_xcellus_replace_row(rawHandle, rawSheetName, rawData *C.char, index int) (failed bool, err *C.char) {
	handle := C.GoString(rawHandle)
	sheetName := C.GoString(rawSheetName)
	data := C.GoString(rawData)

	file, ok := instances[handle]
	if !ok {
		return false, C.CString(invalidHandleErr)
	}
	var dataArr []interface{}
	goErr := json.Unmarshal([]byte(data), &dataArr)

	if goErr != nil {
		return true, errorToCString(goErr)
	}

	sheetPtr, ok := file.Sheet[sheetName]
	if !ok {
		return true, errorToCString("Sheet with provided name does not exist")
	}

	updateRow(sheetPtr, index, dataArr)

	return false, nil
}

func addHeaders(sheetPtr *xlsx.Sheet, wb *worksheet) {
	rowPtr := sheetPtr.AddRow()
	for _, v := range wb.Headers {
		cellPtr := rowPtr.AddCell()
		cellPtr.SetValue(v)
	}
}

func addRows(sheetPtr *xlsx.Sheet, wb *worksheet) {
	for _, r := range wb.Rows {
		rowPtr := sheetPtr.AddRow()
		for _, c := range r {
			cellPtr := rowPtr.AddCell()
			cellPtr.SetValue(c)
		}
	}
}

func findInColumn(sheetPtr *xlsx.Sheet, columnIndex int, value string) int {
	for i := range sheetPtr.Rows {
		if cell := sheetPtr.Rows[i].Cells[columnIndex]; cell != nil {
			if cell.String() == value {
				return i
			}
		}
	}
	return -1
}

func updateRow(sheetPtr *xlsx.Sheet, index int, values []interface{}) {
	row := sheetPtr.Rows[index]
	for i := range values {
		if values[i] == nil {
			continue
		}
		var cell *xlsx.Cell
		if row.Cells[i] == nil {
			cell = row.AddCell()
		} else {
			cell = row.Cells[i]
		}
		cell.SetValue(values[i])
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func errorToCString(err interface{}) *C.char {
	switch err.(type) {
	case error:
		return C.CString(fmt.Sprintf("%s\x00", err))
	case string:
		return C.CString(fmt.Sprintf("%s\x00", err))
	default:
		panic("Invalid type provided to errorTCString")
	}
}
