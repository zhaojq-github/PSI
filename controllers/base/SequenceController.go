package base

import (
	"bytes"
	"encoding/json"
	md "goERP/models"
	"strconv"
	"strings"
)

// SequenceController 登录日志
type SequenceController struct {
	BaseController
}

func (ctl *SequenceController) Post() {
	ctl.URL = "/sequence/"
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

// Get 请求
func (ctl *SequenceController) Get() {

	ctl.URL = "/sequence/"
	ctl.PageName = "序号管理"
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
	ctl.Data["MenuRecordActive"] = "active"

}

// Put update
func (ctl *SequenceController) Put() {
	result := make(map[string]interface{})
	sequence := new(md.Sequence)
	var (
		err error
		id  int64
	)
	if err = ctl.ParseForm(sequence); err == nil {
		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(sequence)).Type().Name()
		if id, err = md.UpdateSequence(sequence, &ctl.User); err == nil {
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

// Edit view
func (ctl *SequenceController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if sequence, err := md.GetSequenceByID(idInt64); err == nil {
				ctl.PageAction = sequence.Name
				ctl.Data["Sequence"] = sequence

			}
		}
	}
	ctl.Data["FormField"] = "form-edit"
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Layout = "base/base.html"
	ctl.TplName = "config/sequence_form.html"
}

// Create view
func (ctl *SequenceController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.Data["FormField"] = "form-create"
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.TplName = "config/sequence_form.html"
}

// Detail view
func (ctl *SequenceController) Detail() {
	//获取信息一样，直接调用Edit
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

// PostCreate post request create sequence
func (ctl *SequenceController) PostCreate() {
	result := make(map[string]interface{})
	sequence := new(md.Sequence)
	var (
		err error
		id  int64
	)
	if err = ctl.ParseForm(sequence); err == nil {
		// 获得struct表名
		// structName := reflect.Indirect(reflect.ValueOf(category)).Type().Name()
		if id, err = md.AddSequence(sequence, &ctl.User); err == nil {
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
func (ctl *SequenceController) Validator() {

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

	if prefix := strings.TrimSpace(ctl.GetString("Prefix")); prefix != "" {
		condAnd["Prefix"] = prefix
	}
	if structName := strings.TrimSpace(ctl.GetString("StructName")); structName != "" {
		condAnd["StructName"] = structName
	}
	if len(condAnd) > 0 {
		cond["and"] = condAnd
	}
	if _, arrs, err := md.GetAllSequence(query, exclude, cond, fields, sortby, order, 0, 2); err == nil {
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

// 获得符合要求的数据
func (ctl *SequenceController) sequenceList(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.Sequence
	paginator, arrs, err := md.GetAllSequence(query, exclude, condMap, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		//使用多线程来处理数据，待修改
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["Name"] = line.Name
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			oneLine["Prefix"] = line.Prefix
			oneLine["Current"] = line.Current
			oneLine["Padding"] = line.Padding
			oneLine["StructName"] = line.StructName
			if line.Company != nil {
				oneLine["Company"] = line.Company.Name
			}
			oneLine["Active"] = line.Active
			oneLine["IsDefault"] = line.IsDefault
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
func (ctl *SequenceController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	name := strings.TrimSpace(ctl.GetString("Name"))
	if name != "" {
		query["Name"] = name
	}
	orderStr := ctl.GetString("order")
	sortStr := ctl.GetString("sort")
	if orderStr != "" && sortStr != "" {
		sortby = append(sortby, sortStr)
		order = append(order, orderStr)
	} else {
		sortby = append(sortby, "Id")
		order = append(order, "desc")
	}
	if result, err := ctl.sequenceList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

// GetList 显示table数据
func (ctl *SequenceController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-sequence"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "config/sequence_list_search.html"

}
