package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GSA/grace-inventory/handler/helpers"
	"github.com/GSA/grace-inventory/handler/inv"
	"github.com/GSA/grace-inventory/handler/spreadsheet"
	"github.com/aws/aws-lambda-go/lambda"
)

var defaultSheets = []string{
	helpers.SheetAccounts,
	helpers.SheetBuckets,
	helpers.SheetGroups,
	helpers.SheetImages,
	helpers.SheetInstances,
	helpers.SheetPolicies,
	helpers.SheetRoles,
	helpers.SheetSecurityGroups,
	helpers.SheetSnapshots,
	helpers.SheetSubnets,
	helpers.SheetUsers,
	helpers.SheetVolumes,
	helpers.SheetVpcs,
	helpers.SheetAddresses,
	helpers.SheetKeyPairs,
	helpers.SheetStacks,
	helpers.SheetAlarms,
	helpers.SheetConfigRules,
	helpers.SheetLoadBalancers,
	helpers.SheetVaults,
	helpers.SheetKeys,
	helpers.SheetDBInstances,
	helpers.SheetDBSnapshots,
	helpers.SheetSecrets,
	helpers.SheetSubscriptions,
	helpers.SheetTopics,
	helpers.SheetParameters,
}

func getSheets() []string {
	v := os.Getenv("sheets")
	if len(v) == 0 {
		return defaultSheets
	}
	sheets := strings.Split(v, ",")

	// prune any references to 'Account' after index zero
	for i := 0; i < len(sheets); i++ {
		if i > 0 && sheets[i] == helpers.SheetAccounts {
			sheets = append(sheets[:i], sheets[i+1:]...)
		}
	}

	// ensure the first element is always 'Accounts'
	if len(sheets) > 0 && sheets[0] != helpers.SheetAccounts {
		sheets = append([]string{helpers.SheetAccounts}, sheets...)
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
