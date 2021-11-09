module github.com/lambda-platform/adminmodule

go 1.15

//replace github.com/lambda-platform/lambda v0.1.8 => ../lambda
//replace github.com/lambda-platform/agent v0.1.9 => ../agent
//replace github.com/lambda-platform/dataform v0.1.1 => ../dataform
//replace github.com/lambda-platform/datagrid v0.1.1 => ../datagrid
//replace github.com/lambda-platform/krud v0.1.0 => ../krud
//

require (
	github.com/labstack/echo/v4 v4.5.0
	github.com/lambda-platform/agent v0.2.2
	github.com/lambda-platform/dataform v0.2.0
	github.com/lambda-platform/datagrid v0.2.0
	github.com/lambda-platform/datasource v0.2.0
	github.com/lambda-platform/krud v0.2.0
	github.com/lambda-platform/lambda v0.2.6
	github.com/lambda-platform/template v0.2.0
)
