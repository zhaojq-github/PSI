package base

import (
	"bytes"
	"encoding/json"
	md "goERP/models"
	"strconv"
	"strings"
)

// CompanyController 公司
type CompanyController struct {
	BaseController
}

// Post request
func (ctl *CompanyController) Post() {
	ctl.URL = "/company/"
	ctl.Data["URL"] = ctl.URL
	action := ctl.Input().Get("action")
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}

// Get request
func (ctl *CompanyController) Get() {
	ctl.URL = "/company/"
	ctl.PageName = "公司管理"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()
	}
	// 标题合成
	b := bytes.Buffer{}
	b.WriteString(ctl.PageName)
	b.WriteString("\\")
	b.WriteString(ctl.PageAction)
	ctl.Data["PageName"] = b.String()
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuCompanyActive"] = "active"
}

// Put update company info
func (ctl *CompanyController) Put() {
	result := make(map[string]interface{})
	company := new(md.Company)
	var (
		err error
		id  int64
	)
	if err = ctl.ParseForm(company); err == nil {

		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(company)).Type().Name()
		if id, err = md.UpdateCompany(company, &ctl.User); err == nil {
			result["code"] = "success"
			result["location"] = ctl.URL + strconv.FormatInt(id, 10) + "?action=detail"
		} else {
			result["code"] = "failed"
			result["message"] = "数据创建失败"
			result["debug"] = err.Error()
		}
	}
	if err != nil {
		result["code"] = "failed"
		result["debug"] = err.Error()
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// PostCreate create company
func (ctl *CompanyController) PostCreate() {
	result := make(map[string]interface{})
	company := new(md.Company)
	var (
		err error
		id  int64
	)
	if err = ctl.ParseForm(company); err == nil {

		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(company)).Type().Name()

		if id, err = md.AddCompany(company, &ctl.User); err == nil {
			result["code"] = "success"
			result["location"] = ctl.URL + strconv.FormatInt(id, 10) + "?action=detail"
		} else {
			result["code"] = "failed"
			result["message"] = "数据创建失败"
			result["debug"] = err.Error()
		}
	} else {
		result["code"] = "failed"
		result["message"] = "请求数据解析失败"
		result["debug"] = err.Error()
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// Edit edit company
func (ctl *CompanyController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if company, err := md.GetCompanyByID(idInt64); err == nil {
				ctl.PageAction = company.Name
				ctl.Data["Company"] = company
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["FormField"] = "form-edit"
	ctl.Layout = "base/base.html"
	ctl.TplName = "user/company_form.html"
}

// Detail detail company
func (ctl *CompanyController) Detail() {
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["FormTreeField"] = "form-tree-edit"
	ctl.Data["Action"] = "detail"
}

// Create create company view
func (ctl *CompanyController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.Data["FormField"] = "form-create"
	ctl.Data["FormTreeField"] = "form-tree-create"
	ctl.TplName = "user/company_form.html"
}

// Validator Validator
func (ctl *CompanyController) Validator() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	condAnd := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)

	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	if parent, err := ctl.GetInt64("Parent"); err == nil {
		query["Parent.Id"] = parent
	}
	if name := strings.TrimSpace(ctl.GetString("Name")); name != "" {
		condAnd["Name"] = name
	}
	if code := strings.TrimSpace(ctl.GetString("Code")); code != "" {
		condAnd["Code"] = code
	}
	if len(condAnd) > 0 {
		cond["and"] = condAnd
	}
	if _, arrs, err := md.GetAllCompany(query, exclude, cond, fields, sortby, order, 0, 2); err == nil {
		if len(arrs) == 1 {
			if arrs[0].ID == recordID {
				result["valid"] = true
			} else {
				result["valid"] = false
			}
		} else {
			result["valid"] = true
		}
	} else {
		result["valid"] = true
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// 获得符合要求的款式数据
func (ctl *CompanyController) companyList(query map[string]interface{}, exclude map[string]interface{}, cond map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.Company
	paginator, arrs, err := md.GetAllCompany(query, exclude, cond, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		//使用多线程来处理数据，待修改
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			oneLine["Code"] = line.Code
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			if line.Parent != nil {
				oneLine["Parent"] = line.Parent.Name
			}
			b := bytes.Buffer{}
			if line.Country != nil {
				b.WriteString(line.Country.Name)
			}
			if line.Province != nil {
				b.WriteString(line.Province.Name)
			}
			if line.City != nil {
				b.WriteString(line.City.Name)
			}
			if line.District != nil {
				b.WriteString(line.District.Name)
			}
			b.WriteString(line.Street)
			oneLine["Address"] = b.String()

			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}

// PostList post request get company list
func (ctl *CompanyController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	condAnd := make(map[string]interface{})
	condOr := make(map[string]interface{})
	filterMap := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
	if ID, err := ctl.GetInt64("Id"); err == nil {
		query["Id"] = ID
	}
	if name := strings.TrimSpace(ctl.GetString("Name")); name != "" {
		condAnd["Name.icontains"] = name
	}
	filter := ctl.GetString("filter")
	if filter != "" {
		json.Unmarshal([]byte(filter), &filterMap)
	}
	if filterName, ok := filterMap["Name"]; ok {
		filterName = strings.TrimSpace(filterName.(string))
		if filterName != "" {
			condAnd["Name.icontains"] = filterName
		}
	}
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	orderStr := ctl.GetString("order")
	sortStr := ctl.GetString("sort")
	if orderStr != "" && sortStr != "" {
		sortby = append(sortby, sortStr)
		order = append(order, orderStr)
	} else {
		sortby = append(sortby, "Id")
		order = append(order, "desc")

	}
	if len(condAnd) > 0 {
		cond["and"] = condAnd
	}
	if len(condOr) > 0 {
		cond["or"] = condOr
	}
	if result, err := ctl.companyList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

// GetList list view
func (ctl *CompanyController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-company"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "user/company_list_search.html"
}
