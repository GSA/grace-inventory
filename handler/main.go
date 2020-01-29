package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GSA/grace-inventory/handler/inv"
	"github.com/GSA/grace-inventory/handler/spreadsheet"
	"github.com/aws/aws-lambda-go/lambda"
)

var defaultSheets = []string{
	inv.SheetAccounts,
	inv.SheetBuckets,
	inv.SheetGroups,
	inv.SheetImages,
	inv.SheetInstances,
	inv.SheetPolicies,
	inv.SheetRoles,
	inv.SheetSecurityGroups,
	inv.SheetSnapshots,
	inv.SheetSubnets,
	inv.SheetUsers,
	inv.SheetVolumes,
	inv.SheetVpcs,
	inv.SheetAddresses,
	inv.SheetKeyPairs,
	inv.SheetStacks,
	inv.SheetAlarms,
	inv.SheetConfigRules,
	inv.SheetLoadBalancers,
	inv.SheetVaults,
	inv.SheetKeys,
	inv.SheetDBInstances,
	inv.SheetDBSnapshots,
	inv.SheetSecrets,
	inv.SheetSubscriptions,
	inv.SheetTopics,
	inv.SheetParameters,
}

func getSheets() []string {
	v := os.Getenv("sheets")
	if len(v) == 0 {
		return defaultSheets
	}
	sheets := strings.Split(v, ",")

	// prune any references to 'Account' after index zero
	for i := 0; i < len(sheets); i++ {
		if i > 0 && sheets[i] == inv.SheetAccounts {
			sheets = append(sheets[:i], sheets[i+1:]...)
		}
	}

	// ensure the first element is always 'Accounts'
	if len(sheets) > 0 && sheets[0] != inv.SheetAccounts {
		sheets = append([]string{inv.SheetAccounts}, sheets...)
	}
	return sheets
}

func createReport() (string, error) {
	filename := fmt.Sprintf("grace_inventory_%s.xlsx", time.Now().Format("2006-01-02-1504"))

	inventory, err := inv.New()
	if err != nil {
		return err.Error(), err
	}

	s := spreadsheet.New(filename)
	sheets := getSheets()
	for _, sheet := range sheets {
		err = s.AddSheet(sheet)
		if err != nil {
			return err.Error(), err
		}
	}

	err = inventory.Run(s)
	if err != nil {
		return err.Error(), err
	}

	return "Report Complete", nil
}

func main() {
	lambda.Start(createReport)
}
