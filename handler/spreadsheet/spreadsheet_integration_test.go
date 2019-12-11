// +build integration

package spreadsheet

import (
	"testing"
)

const itest, itest0, itest1 = "itest", "itest0", "itest1"

// func New(name string) *Spreadsheet
func TestIntegrationNew(t *testing.T) {
	name := itest
	s := New(name)
	result := s.Name
	if result != name {
		t.Fatalf("Name is invalid, expected: %s, got: %s", name, result)
	}
}

// func RegisterSheet(name string, fn SheetFunc)
func TestIntegrationRegisterSheet(t *testing.T) {
	name := itest
	RegisterSheet(name, func() *Sheet { return nil })
	if _, ok := sheetTypes[name]; !ok {
		t.Fatalf("Sheet '%s' not registered", name)
	}
}

// func (ss *Spreadsheet) AddSheet(name string) error
func TestIntegrationAddSheet(t *testing.T) {
	docName := itest0
	sheetName := itest1
	colName := column0
	s := New(docName)
	t.Logf("docName: %v (%T)\n", docName, docName)
	t.Logf("sheetName: %v (%T)\n", sheetName, sheetName)
	t.Logf("s: %#v\n", s)
	err := s.AddSheet(sheetName)
	if err == nil {
		t.Fatal("AddSheet should fail if called before RegisterSheet has been called")
	}
	RegisterSheet(sheetName, func() *Sheet {
		return &Sheet{
			Name: sheetName,
			Columns: []*Column{
				{
					FriendlyName: colName,
					FieldName:    "",
				},
			},
		}
	})
	err = s.AddSheet("not_here")
	if err == nil {
		t.Fatal("AddSheet should fail when referencing an unregistered sheet type")
	}
	if len(s.Sheets) > 0 {
		t.Fatal("a sheet has already been added, Sheets should be empty")
	}
	err = s.AddSheet(sheetName)
	if err != nil {
		t.Fatalf("failed to call AddSheet: %v", err)
	}
	sheet := s.Sheets[0]
	if sheet.Name != sheetName {
		t.Fatalf("sheet 'Name' is invalid, expected: %s, got: %s", sheet.Name, docName)
	}
	column := sheet.Columns[0]
	if column.FriendlyName != colName {
		t.Fatalf("column[0].Name is invalid, expected: %s, got: %s", colName, column.FriendlyName)
	}
}

// func (ss *Spreadsheet) UpdateSheet(name string, payload *Payload)
func testIntegrationUpdateSheet(t *testing.T) {
	docName := itest0
	sheetName := itest1
	colName := column0
	s := New(docName)
	RegisterSheet(sheetName, func() *Sheet {
		return &Sheet{
			Name: sheetName,
			Columns: []*Column{
				{
					FriendlyName: colName,
					FieldName:    "",
				},
				{
					FriendlyName: "name",
					FieldName:    "Name",
				},
				{
					FriendlyName: "value",
					FieldName:    "Value",
				},
			},
		}
	})
	err := s.AddSheet(sheetName)
	if err != nil {
		t.Fatalf("failed to call AddSheet: %v", err)
	}
	objects := []struct {
		Name  string
		Value string
	}{
		{"name0", "value0"},
		{"name1", "value1"},
		{"name2", "value2"},
	}
	var items []interface{}
	for _, o := range objects {
		items = append(items, o)
	}
	s.UpdateSheet(sheetName, &Payload{
		Static: []string{"colval0"},
		Items:  items,
	})
	sheet := s.Sheets[0].sheet
	tests := []struct {
		row      int
		cell     int
		expected string
	}{
		{0, 0, colName},
		{0, 1, "name"},
		{0, 2, "value"},
		{1, 0, "colval0"},
		{1, 1, "name0"},
		{1, 2, "value0"},
		{2, 1, "name1"},
		{2, 2, "value1"},
		{3, 1, "name2"},
		{3, 2, "value2"},
	}
	for _, tt := range tests {
		c := sheet.Cell(tt.row, tt.cell)
		if c.Value != tt.expected {
			t.Fatalf("Cell(%d, %d) invalid, expected: %s, got: %s", tt.row, tt.cell, tt.expected, c.Value)
		}
	}
}

// func (ss *Spreadsheet) Bytes() (*bytes.Reader, error)
func TestIntegrationBytes(t *testing.T) {
	docName := itest0
	s := New(docName)
	_, err := s.Bytes()
	if err == nil {
		t.Fatal("Bytes should fail if no worksheet has been added")
	}
	_, err = s.file.AddSheet(itest)
	if err != nil {
		t.Fatal("failed to add worksheet to spreadsheet")
	}
	r, err := s.Bytes()
	if err != nil {
		t.Fatalf("failed to get bytes: %v", err)
	}
	if r.Len() == 0 {
		t.Fatal("Bytes length invalid, expected: > 0, got: 0")
	}
}

// func (s *Sheet) Update(payload *Payload)
func TestIntegrationUpdate(t *testing.T) {
	docName := itest0
	sheetName := itest1
	colName := column0
	s := New(docName)
	RegisterSheet(sheetName, func() *Sheet {
		return &Sheet{
			Name: sheetName,
			Columns: []*Column{
				{
					FriendlyName: colName,
					FieldName:    "",
				},
				{
					FriendlyName: "name",
					FieldName:    "Name",
				},
				{
					FriendlyName: "value",
					FieldName:    "Value",
				},
			},
		}
	})
	err := s.AddSheet(sheetName)
	if err != nil {
		t.Fatalf("failed to call AddSheet: %v", err)
	}
	objects := []struct {
		Name  string
		Value string
	}{
		{"name0", "value0"},
		{"name1", "value1"},
		{"name2", "value2"},
	}
	var items []interface{}
	for _, o := range objects {
		items = append(items, o)
	}
	s.Sheets[0].Update(&Payload{
		Static: []string{"colval0"},
		Items:  items,
	})
	sheet := s.Sheets[0].sheet
	tests := []struct {
		row      int
		cell     int
		expected string
	}{
		{0, 0, colName},
		{0, 1, "name"},
		{0, 2, "value"},
		{1, 0, "colval0"},
		{1, 1, "name0"},
		{1, 2, "value0"},
		{2, 1, "name1"},
		{2, 2, "value1"},
		{3, 1, "name2"},
		{3, 2, "value2"},
	}
	for _, tt := range tests {
		c := sheet.Cell(tt.row, tt.cell)
		if c.Value != tt.expected {
			t.Fatalf("Cell(%d, %d) invalid, expected: %s, got: %s", tt.row, tt.cell, tt.expected, c.Value)
		}
	}
}
