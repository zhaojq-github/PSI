package address

import (
	"bytes"
	"encoding/json"
	cb "goERP/controllers/base"
	md "goERP/models"
	"strconv"
	"strings"
)

// AddressCityController address city
type AddressCityController struct {
	cb.BaseController
}

// Post request
func (ctl *AddressCityController) Post() {
	ctl.URL = "/address/city/"
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
func (ctl *AddressCityController) Get() {
	ctl.URL = "/address/city/"
	ctl.PageName = "城市管理"
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
	ctl.Data["MenuAddressCityActive"] = "active"
}

// Put put request update city
func (ctl *AddressCityController) Put() {
	result := make(map[string]interface{})
	city := new(md.AddressCity)
	var (
		err error
		id  int64
	)
	if err = ctl.ParseForm(city); err == nil {
		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(city)).Type().Name()
		if id, err = md.UpdateAddressCity(city, &ctl.User); err == nil {
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

// PostCreate post create city
func (ctl *AddressCityController) PostCreate() {
	result := make(map[string]interface{})
	city := new(md.AddressCity)
	var (
		err error
		id  int64
	)
	if err = ctl.ParseForm(city); err == nil {
		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(city)).Type().Name()
		if id, err = md.AddAddressCity(city, &ctl.User); err == nil {
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

// Edit edit view
func (ctl *AddressCityController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if city, err := md.GetAddressCityByID(idInt64); err == nil {
				ctl.PageAction = city.Name
				ctl.Data["City"] = city
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["FormField"] = "form-edit"
	ctl.Layout = "base/base.html"
	ctl.TplName = "address/address_city_form.html"
}

// Detail detail view
func (ctl *AddressCityController) Detail() {
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["FormTreeField"] = "form-tree-edit"
	ctl.Data["Action"] = "detail"
}

// Create create view
func (ctl *AddressCityController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.Data["FormField"] = "form-create"
	ctl.Data["FormTreeField"] = "form-tree-create"
	ctl.TplName = "address/address_city_form.html"
}

// Validator 用于验证
func (ctl *AddressCityController) Validator() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	condAnd := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)

	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	if provinceID, err := ctl.GetInt64("ProvinceID"); err == nil {
		query["Province.Id"] = provinceID
	}
	if name := strings.TrimSpace(ctl.GetString("Name")); name != "" {
		condAnd["Name"] = name
	}
	if _, arrs, err := md.GetAllAddressCity(query, exclude, cond, fields, sortby, order, 0, 2); err == nil {
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
func (ctl *AddressCityController) addressCityList(query map[string]interface{}, exclude map[string]interface{}, cond map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.AddressCity
	paginator, arrs, err := md.GetAllAddressCity(query, exclude, cond, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		//使用多线程来处理数据，待修改
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			oneLine["Province"] = line.Province.Name
			oneLine["Country"] = line.Province.Country.Name
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

// PostList post请求list
func (ctl *AddressCityController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	condAnd := make(map[string]interface{})
	condOr := make(map[string]interface{})
	filterMap := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
	if ProvinceID, err := ctl.GetInt64("ProvinceID"); err == nil {
		query["Province.Id"] = ProvinceID
	}
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
	if len(condAnd) > 0 {
		cond["and"] = condAnd
	}
	if len(condOr) > 0 {
		cond["or"] = condOr
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
	if result, err := ctl.addressCityList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

// GetList list view
func (ctl *AddressCityController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-address-city"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "address/address_city_list_search.html"
}
