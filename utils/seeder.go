package utils

import (
	"encoding/json"
	//"encoding/json"
	"fmt"
	//"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	//puzzleModels "github.com/lambda-platform/lambda/DB/DBSchema/models"
	//analyticModels "github.com/lambda-platform/lambda/lambda/plugins/dataanalytic/models"
	//"os"
	"github.com/lambda-platform/lambda/DB"
	puzzleModels "github.com/lambda-platform/lambda/models"
	"os"
)

func AutoMigrateSeed() {

	if config.Config.Database.Connection == "mssql"{
		DB.DB.AutoMigrate(
			&puzzleModels.VBSchemaMSSQL{},
			&puzzleModels.VBSchemaAdminMSSQL{},
			//&analyticModels.Analytic{},
			//&analyticModels.AnalyticFilter{},
			//&analyticModels.AnalyticRangeFilter{},
			//&analyticModels.AnalyticRowsColumn{},
			//&analyticModels.AnalyticRangeRowColumn{},
			//&analyticModels.AnalyticDateRange{},
		)
	} else {
		DB.DB.AutoMigrate(
			&puzzleModels.VBSchema{},
			&puzzleModels.VBSchemaAdmin{},
			//&analyticModels.Analytic{},
			//&analyticModels.AnalyticFilter{},
			//&analyticModels.AnalyticRangeFilter{},
			//&analyticModels.AnalyticRowsColumn{},
			//&analyticModels.AnalyticRangeRowColumn{},
			//&analyticModels.AnalyticDateRange{},
		)
	}

	if config.Config.App.Seed == "true" {
		var vbs []puzzleModels.VBSchemaAdmin
		DB.DB.Find(&vbs)

		if len(vbs) <= 0 {
			seedData()
		}
	}
}
func seedData() {

	var vbs []puzzleModels.VBSchemaAdmin
	AbsolutePath := AbsolutePath()
	dataFile, err := os.Open(AbsolutePath+"initialData/vb_schemas_admin.json")
	defer dataFile.Close()
	if err != nil {
		fmt.Println("PUZZLE SEED ERROR")
	}
	jsonParser := json.NewDecoder(dataFile)
	err = jsonParser.Decode(&vbs)
	if err != nil {
		fmt.Println(err)
		fmt.Println("PUZZLE SEED DATA ERROR")
	}
	//fmt.Println(len(vbs))

	for _, vb := range vbs {

		DB.DB.Create(&vb)
	}


	var vbs2 []puzzleModels.VBSchema

	dataFile2, err2 := os.Open(AbsolutePath+"initialData/vb_schemas.json")
	defer dataFile2.Close()
	if err2 != nil {
		fmt.Println("PUZZLE SEED ERROR")
	}
	jsonParser2 := json.NewDecoder(dataFile2)
	err = jsonParser2.Decode(&vbs2)
	if err != nil {
		fmt.Println(err)
		fmt.Println("PUZZLE SEED DATA ERROR")
	}
	//fmt.Println(len(vbs))

	for _, vb := range vbs2 {

		DB.DB.Create(&vb)

	}


}
