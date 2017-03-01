package address

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goERP/controllers/base"
	md "goERP/models"
	"strconv"
	"strings"
)

type AddressProvinceController struct {
	base.BaseController
}

func (ctl *AddressProvinceController) Post() {
	ctl.URL = "/address/province/"
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
func (ctl *AddressProvinceController) Get() {
	ctl.URL = "/address/province/"
	ctl.PageName = "省份管理"
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
	ctl.Data["MenuAddressProvinceActive"] = "active"
}

// Put 修改产品款式
func (ctl *AddressProvinceController) Put() {
	result := make(map[string]interface{})

	province := new(md.AddressProvince)
	var (
		err error
		id  int64
	)
	if err = ctl.ParseForm(province); err == nil {
		// 获得struct表名
		fmt.Printf("%+v\n", province)
		// structName := reflect.Indirect(reflect.ValueOf(province)).Type().Name()
		if id, err = md.UpdateAddressProvince(province, &ctl.User); err == nil {
			result["code"] = "success"
			result["location"] = ctl.URL + strconv.FormatInt(id, 10) + "?action=detail"
		} else {
			result["code"] = "failed"
			result["message"] = "数据创建失败"
			result["debug"] = err.Error()
		}
	} else {
		fmt.Println(err)
	}
	if err != nil {
		result["code"] = "failed"
		result["message"] = "解析失败"
		result["debug"] = err.Error()
	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}
func (ctl *AddressProvinceController) PostCreate() {
	result := make(map[string]interface{})
	province := new(md.AddressProvince)
	var (
		err error
		id  int64
	)

	if err = ctl.ParseForm(province); err == nil {
		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(province)).Type().Name()
		if id, err = md.AddAddressProvince(province, &ctl.User); err == nil {
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
func (ctl *AddressProvinceController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if province, err := md.GetAddressProvinceByID(idInt64); err == nil {
				ctl.PageAction = province.Name
				ctl.Data["Province"] = province
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["FormField"] = "form-edit"
	ctl.Layout = "base/base.html"
	ctl.TplName = "address/address_province_form.html"
}
func (ctl *AddressProvinceController) Detail() {
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["FormTreeField"] = "form-tree-edit"
	ctl.Data["Action"] = "detail"
}
func (ctl *AddressProvinceController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.Data["FormField"] = "form-create"
	ctl.Data["FormTreeField"] = "form-tree-create"
	ctl.TplName = "address/address_province_form.html"
}

func (ctl *AddressProvinceController) Validator() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	condAnd := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)

	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	if countryID, err := ctl.GetInt64("CountryID"); err == nil {
		condAnd["Country.Id"] = countryID
	}
	if name := strings.TrimSpace(ctl.GetString("Name")); name != "" {
		condAnd["Name"] = name
	}
	if len(condAnd) > 0 {
		cond["and"] = condAnd
	}
	if _, arrs, err := md.GetAllAddressProvince(query, exclude, cond, fields, sortby, order, 0, 2); err == nil {
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
func (ctl *AddressProvinceController) addressProvinceList(query map[string]interface{}, exclude map[string]interface{}, cond map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.AddressProvince
	paginator, arrs, err := md.GetAllAddressProvince(query, exclude, cond, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		//使用多线程来处理数据，待修改
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			oneLine["Country"] = line.Country.Name
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
func (ctl *AddressProvinceController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	condAnd := make(map[string]interface{})
	condOr := make(map[string]interface{})
	filterMap := make(map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
	if CountryID, err := ctl.GetInt64("CountryID"); err == nil {
		query["Country.Id"] = CountryID
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
	// 对filterMap进行判断
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
	if result, err := ctl.addressProvinceList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

func (ctl *AddressProvinceController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-address-province"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "address/address_province_list_search.html"
}
