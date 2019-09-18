package main

import (
	"fmt"
	"time"

	"github.com/GSA/grace-inventory-lambda/handler/inv"
	"github.com/GSA/grace-inventory-lambda/handler/spreadsheet"
	"github.com/aws/aws-lambda-go/lambda"
)

func createReport() (string, error) {
	filename := fmt.Sprintf("grace_inventory_%s.xlsx", time.Now().Format("2006-01-02-1504"))

	inventory, err := inv.New()
	if err != nil {
		return err.Error(), err
	}

	s := spreadsheet.New(filename)
	sheets := []string{
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
