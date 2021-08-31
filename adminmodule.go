package adminmodule

import (
	"github.com/lambda-platform/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/adminmodule/handlers"
	"github.com/lambda-platform/adminmodule/utils"
	templateUtils "github.com/lambda-platform/template/utils"
	//"github.com/lambda-platform/lambda/lambda/plugins/dataanalytic"
	lambdaUtils "github.com/lambda-platform/lambda/utils"
	"github.com/labstack/echo/v4"
	"html/template"
)

//
func Set(e *echo.Echo, moduleName string, GetGridMODEL func(schema_id string) (interface{}, interface{}, string, string, interface{}, string)) {


	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

	templates := lambdaUtils.GetTemplates(e)
	AbsolutePath := utils.AbsolutePath()
	TemplatePath := templateUtils.AbsolutePath()
	//* REGISTER VIEWS */
	templates["adminmodule.html"] = template.Must(template.ParseFiles(
		TemplatePath+"views/paper.html",
	))
	template.Must(templates["adminmodule.html"].ParseFiles(
		AbsolutePath+"views/adminmodule.html",
	))
	templates["sync_success.html"] = template.Must(template.ParseFiles(
		AbsolutePath+"views/sync_success.html",
	))
	templates["upload_success.html"] = template.Must(template.ParseFiles(
		AbsolutePath+"views/upload_success.html",
	))

	/*ROUTES */
	e.GET("/build-me", handlers.BuildMe, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g := e.Group("/lambda")

	//g.GET("/puzzle", handlers.Index, agentMW.IsLoggedInCookie)
	g.GET("/puzzle", handlers.Index, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.GET("/upload", handlers.UploadDBSCHEMA, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.GET("/async", handlers.ASyncFromCloud, agentMW.IsLoggedInCookie, agentMW.IsAdmin)

	//Puzzle
	g.GET("/puzzle/schema/:type", handlers.GetVB, agentMW.IsLoggedInCookie)
	g.GET("/puzzle/schema/:type/:id", handlers.GetVB, agentMW.IsLoggedInCookie)
	g.GET("/puzzle/schema-public/:type/:id", handlers.GetVB)
	g.GET("/puzzle/schema/:type/:id/:condition", handlers.GetVB, agentMW.IsLoggedInCookie)

	//VB SCHEMA
	g.GET("/puzzle/table-schema/:table", handlers.GetTableSchema, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.POST("/puzzle/schema/:type", handlers.SaveVB(moduleName), agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.POST("/puzzle/schema/:type/:id", handlers.SaveVB(moduleName), agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.DELETE("/puzzle/delete/vb_schemas/:type/:id", handlers.DeleteVB, agentMW.IsLoggedInCookie, agentMW.IsAdmin)

	//GRID
	g.POST("/puzzle/grid/:action/:schemaId", handlers.GridVB(GetGridMODEL), agentMW.IsLoggedInCookie)

	//Get From Options
	g.POST("/puzzle/get_options", handlers.GetOptions, agentMW.IsLoggedInCookie)
	g.POST("/puzzle/get_options-public", handlers.GetOptions)

	//Roles
	g.GET("/puzzle/roles-menus", handlers.GetRolesMenus, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.GET("/puzzle/get-krud-fields/:id", handlers.GetKrudFields, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.POST("/puzzle/save-role", handlers.SaveRole, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.POST("/puzzle/roles/create", handlers.CreateRole, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.POST("/puzzle/roles/store/:id", handlers.UpdateRole, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	g.DELETE("/puzzle/roles/destroy/:id", handlers.DeleteRole, agentMW.IsLoggedInCookie, agentMW.IsAdmin)

}
