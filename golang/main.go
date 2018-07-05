package main

/*
#include <stdlib.h>
*/
import "C"

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"

    "github.com/tealeg/xlsx"
)

func main() {}

type worksheet struct {
    Title   string          `json:"title"`
    Headers []string        `json:"headers"`
    Rows    [][]interface{} `json:"rows"`
}

//export go_xcellus_process
func go_xcellus_process(data *C.char) (failed bool, err *C.char, len int, buf *C.char) {
    synthesizedData := C.GoString(data)
    var wb []*worksheet
    goErr := json.Unmarshal([]byte(synthesizedData), &wb)

    if goErr != nil {
        return true, C.CString(fmt.Sprintf("%s\x00", goErr)), 0, nil
    }

    file := xlsx.NewFile()
    var sheetPtr *xlsx.Sheet
    var rowPtr *xlsx.Row
    var cellPtr *xlsx.Cell

    for _, sheet := range wb {
        sheetPtr, goErr = file.AddSheet(sheet.Title)
        if goErr != nil {
            return true, C.CString(fmt.Sprintf("%s\x00", goErr)), 0, nil
        }
        rowPtr = sheetPtr.AddRow()
        for _, v := range sheet.Headers {
            cellPtr = rowPtr.AddCell()
            cellPtr.SetValue(v)
        }
        for _, r := range sheet.Rows {
            rowPtr = sheetPtr.AddRow()
            for _, c := range r {
                cellPtr = rowPtr.AddCell()
                cellPtr.SetValue(c)
            }
        }
    }

    tmpBuffer := &bytes.Buffer{}
    writer := bufio.NewWriter(tmpBuffer)
    goErr = file.Write(writer)

    if goErr != nil {
        return true, C.CString(fmt.Sprintf("%s\x00", goErr)), 0, nil
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
