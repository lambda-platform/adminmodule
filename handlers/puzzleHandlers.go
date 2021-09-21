package handlers

import (
	"encoding/json"
	"database/sql"
	"fmt"
	agentUtils "github.com/lambda-platform/agent/utils"
	"github.com/lambda-platform/datagrid"
	"github.com/lambda-platform/dataform"
	"archive/zip"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/models"
	"github.com/lambda-platform/lambda/DBSchema"
	"os/exec"

	"github.com/lambda-platform/lambda/config"
	"github.com/labstack/echo/v4"
	"strings"
	"os"
	"github.com/lambda-platform/datasource"
	"github.com/lambda-platform/lambda/utils"
	"net/http"
	"regexp"
	"strconv"
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"io/ioutil"
	krudModels "github.com/lambda-platform/krud/models"
)

type vb_schema struct {
	ID         int        `gorm:"column:id;primary_key" json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

func Index(c echo.Context) error {
	dbSchema := models.DBSCHEMA{}

	if(config.LambdaConfig.SchemaLoadMode == "auto"){
		dbSchema = DBSchema.GetDBSchema()
	} else {
		schemaFile, err := os.Open("app/models/db_schema.json")
		defer schemaFile.Close()
		if err != nil{
			fmt.Println("schema FILE NOT FOUND")
		}
		dbSchema = models.DBSCHEMA{}
		jsonParser := json.NewDecoder(schemaFile)
		jsonParser.Decode(&dbSchema)
	}



	userRoles := []models.UserRoles{}


	DB.DB.Find(&userRoles)

	//gridList, err := models.VBSchemas(qm.Where("type = ?", "grid")).All(context.Background(), DB)
	//dieIF(err)

	User := agentUtils.AuthUserObject(c)

	//csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	csrfToken := ""
	return c.Render(http.StatusOK, "adminmodule.html", map[string]interface{}{
		"title":                     config.LambdaConfig.Title,
		"lambda_config": config.LambdaConfig,
		"favicon":                     config.LambdaConfig.Favicon,
		"app_logo":                     config.LambdaConfig.Logo,
		"app_text":                     "СИСТЕМИЙН УДИРДЛАГА",
		"dbSchema":                  dbSchema,
		"User":                      User,
		"user_fields":               config.LambdaConfig.UserDataFields,
		"user_roles":               userRoles,
		"data_form_custom_elements": config.LambdaConfig.DataFormCustomElements,
		"mix":                       utils.Mix,
		"csrfToken":                       csrfToken,
	})

}


func GetTableSchema(c echo.Context) error {
	table := c.Param("table")
	tableMetas := DBSchema.TableMetas(table)
	return c.JSON(http.StatusOK, tableMetas)

}
func UploadDBSCHEMA(c echo.Context) error {

	DBSchema.GenerateSchemaForCloud()


	url := "https://lambda.cloud.mn/console/upload/"+config.LambdaConfig.ProjectKey
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("app/models/db_schema.json")
	defer file.Close()
	part1,
	errFile1 := writer.CreateFormFile("file",filepath.Base("app/models/db_schema.json"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	}


	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})

	}



	return c.Render(http.StatusOK, "upload_success.html", map[string]interface{}{
		"status": true,
		"body": body,
	})
}
func ASyncFromCloud(c echo.Context) error {




	url := "https://lambda.cloud.mn/console/project-data/"+config.LambdaConfig.ProjectKey

	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	}
	fmt.Println("=============== HERE 1")

	data := CloudData{}
	json.Unmarshal(body, &data)
	fmt.Println("=============== HERE 2")
	DSVBS := []models.VBSchema{}
	FormVbs := []models.VBSchema{}
	GridVbs := []models.VBSchema{}
	MenuVbs := []models.VBSchema{}
	cruds := []krudModels.Krud{}
	FormSchemasJSON, _ := json.Marshal(data.FormSchemas)
	GridSchemasJSON, _ := json.Marshal(data.GridSchemas)
	MenuSchemasJSON, _ := json.Marshal(data.MenuSchemas)
	KrudJSON, _ := json.Marshal(data.Cruds)
	json.Unmarshal([]byte(FormSchemasJSON), &FormVbs)
	json.Unmarshal([]byte(GridSchemasJSON), &GridVbs)
	json.Unmarshal([]byte(MenuSchemasJSON), &MenuVbs)
	json.Unmarshal([]byte(KrudJSON), &cruds)
	fmt.Println("=============== HERE 3")
	DB.DB.Where("type = ?", "datasource").Find(&DSVBS)
	fmt.Println("=============== HERE 3.1")
	DB.DB.Exec("TRUNCATE krud")
	DB.DB.Exec("TRUNCATE vb_schemas")
	fmt.Println("=============== HERE 4")
	for _, vb := range FormVbs {
		DB.DB.Create(&vb)
	}
	for _, vb := range GridVbs {
		DB.DB.Create(&vb)
	}
	for _, vb := range MenuVbs {
		DB.DB.Create(&vb)
	}
	for _, vb := range DSVBS {
		vb.ID = 0;
		DB.DB.Create(&vb)

	}
	for _, crud := range cruds {

		DB.DB.Create(&crud)
	}
	fmt.Println("=============== HERE 5")
	var downloadError error = DownloadGeneratedCodes()
	fmt.Println("=============== HERE 6")
	if(downloadError != nil){
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": false,
			"msg": downloadError,
		})
	} else {

		var unzip error = UnZipLambdaCodes()


		if(unzip != nil){
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": false,
				"msg": unzip,
			})
		} else {
			ReBuild()
			return c.Render(http.StatusOK, "sync_success.html", map[string]interface{}{
				"status": true,
			})
		}
	}



}
func ReBuild() {
	//bytes, err1 := exec.Command("killall", "lambda-starter").Output()
	//output := string(bytes)
	//fmt.Println(output)
	//if err1 != nil {
	//	fmt.Println(err1)
	//}

	dir, _ := os.Getwd()
	fmt.Println(dir)
	out, err := exec.Command("/bin/sh",  dir+"/script-build.sh").Output()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("The date is %s", out)

}

func DownloadGeneratedCodes() error{
	url := "https://lambda.cloud.mn/console/get-codes/"+config.LambdaConfig.ProjectKey

	resp, err := http.Get(url)
	if err != nil {
		return err
	}


	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("file not found error")
	}

	// Create the file
	out, err := os.Create("lambda.zip")
	if err != nil {

		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)


	return err
}
func UnZipLambdaCodes() error{
	var dest string = "lambda"
	var src string = "lambda.zip"
	if(!utils.FileExists(src)){
		return errors.New("Lambda file Not found")
	} else {
		formPatch :="lambda/models/form/"
		gridPatch := "lambda/models/grid/"
		if _, err := os.Stat(formPatch); os.IsNotExist(err) {
			os.MkdirAll("lambda/models/form", 0755)
			os.MkdirAll(formPatch, 0755)
			os.MkdirAll("lambda/models/form/validationCaller/", 0755)
			os.MkdirAll("lambda/models/form/validations/", 0755)
			os.MkdirAll("lambda/models/form/caller/", 0755)
		} else {
			os.MkdirAll("lambda/models/form", 0755)
			os.RemoveAll(formPatch)
			os.MkdirAll(formPatch, 0755)
			os.MkdirAll("lambda/models/form/validationCaller/", 0755)
			os.MkdirAll("lambda/models/form/validations/", 0755)
			os.MkdirAll("lambda/models/form/caller/", 0755)
		}
		if _, err := os.Stat(gridPatch); os.IsNotExist(err) {
			os.MkdirAll("lambda/models/grid", 0755)
			os.MkdirAll(gridPatch, 0755)
			os.MkdirAll("lambda/models/grid/caller", 0755)
		} else {
			os.MkdirAll("lambda/models/grid", 0755)
			os.RemoveAll(gridPatch)
			os.MkdirAll(gridPatch, 0755)
			os.MkdirAll("lambda/models/grid/caller", 0755)
		}

		graphqlPatch :=  "lambda/graph"
		graphqlGeneratedPatch :=  "lambda/graph/generated"
		modelsPatch :=  "lambda/graph/models"
		schemaPatch :=  "lambda/graph/schemas"
		resolversPatch :=  "lambda/graph/resolvers"
		schemaCommonPatch :=  "lambda/graph/schemas-common"
		if _, err := os.Stat(modelsPatch); os.IsNotExist(err) {

			os.MkdirAll(graphqlPatch, 0755)
			os.MkdirAll(graphqlGeneratedPatch, 0755)
			os.MkdirAll(modelsPatch, 0755)
			os.MkdirAll(schemaPatch, 0755)
			os.MkdirAll(resolversPatch, 0755)
			os.MkdirAll(schemaCommonPatch, 0755)

		} else {

			os.RemoveAll(graphqlPatch)
			os.RemoveAll(graphqlGeneratedPatch)
			os.RemoveAll(modelsPatch)
			os.RemoveAll(schemaPatch)
			os.RemoveAll(resolversPatch)
			os.RemoveAll(schemaCommonPatch)
			os.MkdirAll(graphqlPatch, 0755)
			os.MkdirAll(graphqlGeneratedPatch, 0755)
			os.MkdirAll(modelsPatch, 0755)
			os.MkdirAll(schemaPatch, 0755)
			os.MkdirAll(resolversPatch, 0755)
			os.MkdirAll(schemaCommonPatch, 0755)
		}

		var filenames []string

		r, err := zip.OpenReader(src)
		if err != nil {
			return  err
		}
		defer r.Close()

		for _, f := range r.File {

			// Store filename/path for returning and using later on
			fpath := filepath.Join(dest, f.Name)

			// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
			if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
				return fmt.Errorf("%s: illegal file path", fpath)
			}

			filenames = append(filenames, fpath)

			if f.FileInfo().IsDir() {
				// Make Folder
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return  err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			rc.Close()

			if err != nil {
				return  err
			}
		}
		e := os.Remove(src)
		if e != nil {
			return e
		}
		return nil
	}







}
type CloudData struct {
	Cruds []struct {
		Form       int    `json:"form"`
		Grid       int    `json:"grid"`
		ID         int    `json:"id"`
		ProjectsID int    `json:"projects_id"`
		Template   string `json:"template"`
		Title      string `json:"title"`
	} `json:"cruds"`
	GridSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"form-schemas"`
	FormSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"grid-schemas"`
	MenuSchemas []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ProjectsID int    `json:"projects_id"`
		Schema     string `json:"schema"`
		Type       string `json:"type"`
	} `json:"menu-schemas"`
}
func GridVB(GetGridMODEL func(schema_id string) (interface{}, interface{}, string, string, interface{}, string)) echo.HandlerFunc {
	return func(c echo.Context) error {
		schemaId := c.Param("schemaId")
		action := c.Param("action")
		id := c.Param("id")

		return datagrid.Exec(c, schemaId, action, id, GetGridMODEL)
	}
}
func GetOptions(c echo.Context) error {

	r := new(dataform.Relations)
	if err := c.Bind(r); err != nil {

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": false,
			"error": err.Error(),
		})
	}
	optionsData := map[string][]map[string]interface{}{}

	var DB_ *sql.DB
	DB_ = DB.DB.DB()
	for table, relation := range r.Relations {
		data := dataform.OptionsData(DB_, relation, c)
		optionsData[table] = data

	}
	return c.JSON(http.StatusOK, optionsData)

}
func GetVB(c echo.Context) error {

	type_ := c.Param("type")
	id := c.Param("id")
	condition := c.Param("condition")

	if id != "" {

		match, _ := regexp.MatchString("_", id)

		if(match){
			VBSchema := models.VBSchemaAdmin{}

			DB.DB.Where("id = ?", id).First(&VBSchema)

			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": true,
				"data":   VBSchema,
			})
		} else {

			VBSchema := models.VBSchema{}

			DB.DB.Where("id = ?", id).First(&VBSchema)

			if type_ == "form"{

				if condition != ""{
					if condition != "builder"{
						return dataform.SetCondition(condition, c, VBSchema)
					}
				}
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": true,
				"data":   VBSchema,
			})
		}





	} else {

		VBSchemas := []models.VBSchemaList{}

		DB.DB.Select("id, name, type, created_at, updated_at").Where("type = ?", type_).Find(&VBSchemas)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": true,
			"data":   VBSchemas,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": false,
	})

}
func SaveVB(modelName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type_ := c.Param("type")
		id := c.Param("id")
		//condition := c.Param("condition")

		vbs := new(vb_schema)
		if err := c.Bind(vbs); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": false,
				"error": err.Error(),
			})
		}

		if id != "" {
			id_, _ := strconv.ParseUint(id, 0, 64)

			vb := models.VBSchema{}

			DB.DB.Where("id = ?", id_).First(&vb)

			vb.Name = vbs.Name
			vb.Schema = vbs.Schema
			//_, err := vb.Update(context.Background(), DB, boil.Infer())

			BeforeSave(id_, type_)

			err := DB.DB.Save(&vb).Error

			if type_ == "form" {
				//WriteModelData(id_)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModel(modelName)
				//WriteGridModelById(modelName, vb.ID)
			}

			if err != nil {

				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"status": false,
					"error": err.Error(),
				})
			} else {

				error := AfterSave(vb, type_)

				if(error != nil){
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": false,
						"error":error.Error(),
					})
				} else {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": true,
					})
				}
			}

		} else {
			vb := models.VBSchema{
				Name:   vbs.Name,
				Schema: vbs.Schema,
				Type:   type_,
				ID:0,
			}

			//err := vb.Insert(context.Background(), DB, boil.Infer())

			DB.DB.NewRecord(vb) // => returns `true` as primary key is blank

			err := DB.DB.Create(&vb).Error

			if type_ == "form" {
				//WriteModelData(vb.ID)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModelById(modelName, vb.ID)
				//WriteGridModel(modelName)
			}



			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"status": "false",
				})
			} else {
				error := AfterSave(vb, type_)

				if(error != nil){
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": false,
						"error":error.Error(),
					})
				} else {
					return c.JSON(http.StatusOK, map[string]interface{}{
						"status": true,
					})
				}

			}

		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": false,
		})
	}
}

func DeleteVB(c echo.Context) error {

	type_ := c.Param("type")
	id := c.Param("id")
	//condition := c.Param("condition")

	vbs := new(vb_schema)
	id_, _ := strconv.ParseUint(id, 0, 64)

	BeforeDelete(id_, type_)

	err := DB.DB.Where("id = ?", id).Where("type = ?", type_).Delete(&vbs).Error

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "true",
		})
	}



}

func BeforeDelete(id uint64, type_ string){

	if type_ == "datasource"{
		vb := models.VBSchema{}

		DB.DB.Where("id = ?", id).First(&vb)

		datasource.DeleteView("ds_"+vb.Name)
	}

}
func BeforeSave(id uint64, type_ string){

	if type_ == "datasource"{
		vb := models.VBSchema{}

		DB.DB.Where("id = ?", id).First(&vb)

		datasource.DeleteView("ds_"+vb.Name)
	}

}
func AfterSave(vb models.VBSchema, type_ string) error{

	if type_ == "datasource"{
		return datasource.CreateView(vb.Name, vb.Schema)
	}

	return nil

}


