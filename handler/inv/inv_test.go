package inv

//"github.com/GSA/grace-inventory/handler/spreadsheet"
//awstest "github.com/gruntwork-io/terratest/modules/aws"

//var genericRegions = []string{"us-east-1", "us-west-2"}

// New(defaultRegion string, regions []string, mgmtAccount string, bucketID string, kmsKeyID string) (*Inv, error)
/*
func TestNew(t *testing.T) {
	i, err := New()
	if err != nil {
		t.Fatalf("failed to create a new session: %v", err)
	}
	if i.defaultRegion != genericRegions[0] {
		t.Fatalf("defaultRegion invalid, expected: %s, got: %s", genericRegions[0], i.defaultRegion)
	}
	if len(i.regions) != len(genericRegions) {
		t.Fatalf("regions length invalid, expected: %d, got: %d", len(genericRegions), len(i.regions))
	}
}
*/
// func (inv *Inv) Run(s *spreadsheet.Spreadsheet) error
/*
func TestRun(t *testing.T) {
	currUser := awstest.GetIamCurrentUserName(t)
	i, err := New(genericRegions[0], genericRegions, currUser, bucketId, kmsKeyId)
	if err != nil {
		t.Fatalf("failed to create a new session: %v", err)
	}
	s := spreadsheet.New("test")
	s.AddSheet(SHEET_BUCKETS)
	err = i.Run(s)
	if err != nil {
		t.Fatalf("failed to execute Run: %v", err)
	}
}
*/
