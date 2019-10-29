package spreadsheet

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/tealeg/xlsx"
)

// SheetFunc ... used by RegisterSheet to register
// sheets allow to be used by the AWS report
type SheetFunc func() *Sheet

var sheetTypes map[string]SheetFunc

// Payload ... used by Update and UpdateSheet to populate
// a sheet with particular datasets. Static is prepended
// to every row created by Items. Items should be a slice
// of objects
type Payload struct {
	Static []string
	Items  []interface{}
}

func (p *Payload) String() string {
	var sb strings.Builder
	_, err := sb.WriteString(fmt.Sprintf("Static: %v,\n Items: [", p.Static))
	if err != nil {
		fmt.Printf("Warning: Error writing string: %v", err)
	}
	for _, i := range p.Items {
		_, err = sb.WriteString(fmt.Sprintf(`%#v`, i))
		if err != nil {
			fmt.Printf("Warning: Error writing string: %v", err)
		}
		_, err = sb.WriteString(",\n")
		if err != nil {
			fmt.Printf("Warning: Error writing string: %v", err)
		}
	}
	_, err = sb.WriteString("]\n")
	if err != nil {
		fmt.Printf("Warning: Error writing string: %v", err)
	}
	return sb.String()
}

// Column ... used to describe a column on a sheet
// if FieldName is empty, the column is considered
// to be static
type Column struct {
	FriendlyName string
	FieldName    string
}

// Spreadsheet ... holds the desired filename and
// all sheets created by calling 'AddSheet'
type Spreadsheet struct {
	file   *xlsx.File
	Name   string
	Sheets []*Sheet
}

// New ... returns a *Spreadsheet, and sets the filename
// to the provided 'name'
func New(name string) *Spreadsheet {
	return &Spreadsheet{
		file: xlsx.NewFile(),
		Name: name,
	}
}

// RegisterSheet ... registers a SheetFunc with the given
// 'name'. A SheetFunc must be registered before calling
// 'AddSheet' to add sheets.
func RegisterSheet(name string, fn SheetFunc) {
	if sheetTypes == nil {
		sheetTypes = make(map[string]SheetFunc)
	}
	sheetTypes[name] = fn
}

// AddSheet ... creates a new sheet using the matching SheetFunc
// given the provided 'name'. Then initializes the sheet by creating
// the header row and adding the column names
func (ss *Spreadsheet) AddSheet(name string) error {
	if sheetTypes == nil {
		return errors.New("zero Sheet Types have been registered")
	}
	if fn, ok := sheetTypes[name]; ok {
		var err error
		s := fn()
		// add sheet to underlying file with friendlyName
		s.sheet, err = ss.file.AddSheet(s.Name)
		// update our local sheet's name with the internal name
		s.Name = name
		if err != nil {
			return err
		}
		row := s.sheet.AddRow()
		for _, c := range s.Columns {
			cell := row.AddCell()
			cell.Value = c.FriendlyName
		}
		ss.Sheets = append(ss.Sheets, s)
		return nil
	}
	return fmt.Errorf("%s is not a registered Sheet Type", name)
}

// UpdateSheet ... finds the sheet matching the given 'name', then calls
// Update passing the payload provided
func (ss *Spreadsheet) UpdateSheet(name string, payload *Payload) {
	for _, s := range ss.Sheets {
		if s.Name == name {
			s.Update(payload)
			return
		}
	}
}

// Bytes ... creates a bytes.Buffer, saves the underlying xlsx.File
// to the buffer, then returns the bytes wrapped in a bytes.Reader
func (ss *Spreadsheet) Bytes() (*bytes.Reader, error) {
	// Hopefully this works... Create a buffer,
	// write the document to it, then wrap the bytes in a bytes.Reader
	buf := bytes.Buffer{}
	err := ss.file.Write(&buf)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// Sheet ... holds a pointer to the underlying xlsx.Sheet, the sheet name,
// and all of the columns returned by the SheetFunc
type Sheet struct {
	sheet   *xlsx.Sheet
	Name    string
	Columns []*Column
}

// Update ... Enumerates over the provided array, adding a new row for each
// element and prepending each row with the StaticValues
func (s *Sheet) Update(payload *Payload) {
	if payload == nil {
		return
	}
	for _, obj := range payload.Items {
		row := s.sheet.AddRow()
		for _, s := range payload.Static {
			cell := row.AddCell()
			cell.Value = s
		}
		for _, c := range s.Columns {
			if c.FieldName == "" {
				continue
			}

			cell := row.AddCell()
			cell.Value = ""

			val := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(c.FieldName)
			// handle nil here instead of inside setCell
			if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
				continue
			}
			s.setCell(cell, val.Interface())
		}
	}
}

// getTagName ... loops over tags looking for a Key that matches Name and returns the Value
func getTagName(tags []*ec2.Tag) string {
	for _, t := range tags {
		if t != nil && aws.StringValue(t.Key) == "Name" {
			return aws.StringValue(t.Value)
		}
	}
	return ""
}

// nolint: gocyclo
// setCell ... sets the value of a cell, after converting it from interface{}
func (s *Sheet) setCell(cell *xlsx.Cell, val interface{}) {
	switch v := val.(type) {
	case nil:
		return
	case *string:
		cell.Value = aws.StringValue(v)
	case *bool:
		cell.SetBool(aws.BoolValue(v))
	case *int:
		cell.SetInt(aws.IntValue(v))
	case *int64:
		cell.SetInt64(aws.Int64Value(v))
	case *float64:
		cell.SetFloat(aws.Float64Value(v))
	case *ec2.InstanceState:
		cell.Value = aws.StringValue(v.Name)
	case int:
		cell.SetInt(v)
	case int64:
		cell.SetInt64(v)
	case float64:
		cell.SetFloat(v)
	case string:
		cell.Value = v
	case time.Time:
		cell.SetDateTime(v)
	case []*ec2.Tag:
		cell.Value = getTagName(v)
	case *time.Time:
		cell.SetDateTime(aws.TimeValue(v))
	}
}
