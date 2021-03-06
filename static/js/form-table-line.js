// 删除某一行
var formTreeLineActionRemove = function() {
    $(".form-tree-line-action-remove").on('click', function(e) {
        e.preventDefault();
        var adjustTr = e.currentTarget;
        for (var i = 0; i < 5; i++) {
            if (adjustTr.nodeName != "TR") {
                adjustTr = adjustTr.parentNode;
            } else {
                $(adjustTr).addClass("danger");
            }
        }
    })
};
formTreeLineActionRemove();
displayTable("#form-table-sale-order-line", "/sale/order/line", [
    { title: "全选", field: 'ID', checkbox: true, align: "center", valign: "middle" },
    { title: "订单明细号", field: 'Name', align: "left", sortable: true, order: "desc", valign: "middle" },
    {
        title: "产品",
        field: 'Product',
        align: "left",
        sortable: true,
        order: "desc",
        valign: "middle",
        cellStyle: { css: { "min-width": "200px" } },
        formatter: function cellStyle(value, row, index) {
            var html = '';
            if (row.Product) {
                html = "<p class='p-form-line-control'>" + "[" + row.Product.defaultCode + "]" + row.Product.name + "<a class='pull-right' target='_blank' href='/product/product/" + row.Product.id + "?action=detail'><i class='fa fa-external-link'></i></a></p>";
                html += '<select   name="ProductProduct-' + row.id + '" id="ProductProduct-' + row.id + '" class="form-control select-sale-order-product-product">';
                html += '<option value="' + row.Product.id + '"  selected="selected">' + '[' + row.Product.defaultCode + ']' + row.Product.name + '</option>'
                html += '</select>';
            }
            return html;
        }
    },
    {
        title: "产品编码",
        field: 'ProductCode',
        align: "left",
        sortable: true,
        order: "desc",
        valign: "middle",
        formatter: function cellStyle(value, row, index) {
            var html = '';
            if (row.ProductCode) {
                html = "<p class='p-form-line-control'>" + row.ProductCode + "</p>";
            }
            return html;
        }
    },
    {
        title: "产品名称",
        field: 'ProductName',
        align: "left",
        sortable: true,
        order: "desc",
        valign: "middle",
        formatter: function cellStyle(value, row, index) {
            var html = '';
            if (row.ProductName) {
                html = "<p class='p-form-line-control'>" + row.ProductName + "</p>";
            }
            return html;
        }
    },
    {
        title: "第一单位数量",
        field: 'FirstSaleQty',
        align: "center",
        sortable: true,
        order: "desc",
        valign: "middle",
        cellStyle: { css: { "max-width": "40px" } },
        formatter: function cellStyle(value, row, index) {
            var html = '';
            var firstSaleQty = row.FirstSaleQty || 1;
            var firstUomStep = row.FirstUomStep || 1;
            var firstUomName = row.FirstUomName;
            var firstUomNameStr = "";
            if (firstUomName != undefined || firstUomName != "") {
                firstUomNameStr = "<span class='pull-right'>(" + firstUomName + ")</span>";
            }
            html = "<p class='p-form-line-control'>" + firstSaleQty + firstUomNameStr + "</p>";
            html += '<input class="form-control " id="FirstSaleQty-' + row.id + ' name="FirstSaleQty-' + row.id + '" type="number" min="0" step="' + firstUomStep + '" value="' + firstSaleQty + '"/>';
            return html;
        }
    },
    {
        title: "第二单位数量",
        field: 'SecondSaleQty',
        align: "center",
        sortable: true,
        order: "desc",
        valign: "middle",
        cellStyle: { css: { "max-width": "40px" } },
        formatter: function cellStyle(value, row, index) {
            var html = '';
            var secondSaleQty = row.SecondSaleQty || 0;
            var secondUomStep = row.SecondUomStep || 1;
            var secondUomName = row.SecondUomName;
            var secondUomNameStr = "";
            if (secondUomName != undefined || secondUomName != "") {
                secondUomNameStr = "<span class='pull-right'>(" + secondUomName + ")</span>";
            }
            html = "<p class='p-form-line-control'>" + secondSaleQty + secondUomNameStr + "</p>";
            html += '<input class="form-control " id="SecondSaleQty-' + row.id + ' name="SecondSaleQty-' + row.id + '" type="number" min="0" step="' + secondUomStep + '" value="' + secondSaleQty + '"/>';
            return html;
        }
    },
    {
        title: "单价",
        field: 'PriceUnit',
        align: "left",
        sortable: true,
        order: "desc",
        valign: "middle",
        formatter: function cellStyle(value, row, index) {
            var html = '';
            if (row.PriceUnit) {
                html = "<p class='p-form-line-control'>" + row.PriceUnit + "</p>";
            }
            return html;
        }
    },
    {
        title: "小计",
        field: 'Total',
        align: "left",
        sortable: true,
        order: "desc",
        valign: "middle",
        formatter: function cellStyle(value, row, index) {
            var html = '';
            if (row.Total) {
                html = "<p class='p-form-line-control'>" + row.Total + "</p>";
            }
            return html;
        }
    },
    {
        title: "操作",
        align: "center",
        field: 'action',
        formatter: function cellStyle(value, row, index) {
            var html = "";
            var url = "/sale/order/line/";
            html += "<a class='form-tree-line-action-remove' href='#'><i class='fa fa-trash-o'></i></a>";
            return html;
        }
    }
], {
    onPostBody: function() {
        select2AjaxData(".select-sale-order-product-product", '/product/product/', {
            changeFunction: function(event) {},
            formatRepo: function(repo) {
                'use strict';
                if (repo.loading) { return repo.text; }
                var name = repo.name || "";
                if (repo.DefaultCode != undefined) {
                    name += '[' + repo.DefaultCode + ']';
                }
                if (repo.Name != undefined) {
                    name += repo.Name;
                }
                return "<p>" + name + "</p>";
            },
            formatRepoSelection: function(repo) {
                'use strict';
                var name = repo.name || "";
                if (repo.DefaultCode != undefined) {
                    name += '[' + repo.DefaultCode + ']';
                }
                if (repo.Name != undefined) {
                    name += repo.Name;
                }
                return "<p>" + name + "</p>";
            }
        });
    },
    queryParams: function(params) {
        var xsrf = $("input[name ='_xsrf']");
        if (xsrf.length > 0) {
            params._xsrf = xsrf[0].value;
        }
        params.action = 'table';

        var saleOrderId = $("input[name ='recordID']");
        if (saleOrderId.length > 0) {
            params.saleOrderId = parseInt(saleOrderId[0].value);
        } else {
            params.saleOrderId = 0;
        }
        return params;
    }
});
// 增加一行销售订单明细
$("#add-one-sale-order-line").on("click", function(e) {
    $("#form-table-sale-order-line").bootstrapTable('prepend', [{
        FirstSaleQty: 1,
        ID: null,
        Name: "",
        PriceUnit: 0,
        Product: {
            id: null,
            name: "",
            defaultCode: ""
        },
        ProductCode: "",
        ProductName: "",
        SecondSaleQty: 0,
        Total: 0,
        id: null
    }]);
});
displayTable("#product-template-attribute-table", "/product/attributeline/", [
    { title: "全选", field: 'ID', checkbox: true, align: "center", valign: "middle" },
    {
        title: "属性",
        field: 'Attribute',
        align: "left",
        sortable: true,
        order: "desc",
        valign: "middle",
        formatter: function(value, row, index) {
            var html = "<select class='form-table-product-template-attributte'>";
            if (row.Attribute && row.Attribute.id) {
                html += +"<option value=" + row.Attribute.id + ">" + row.Attribute.Name + "</option>";
            }
            html += "</select>";
            return html;
        }
    },
    {
        title: "属性值",
        field: 'AttributeValues',
        align: "left",
        sortable: true,
        order: "desc",
        valign: "middle",
        cellStyle: { css: { "min-width": "200px" } },
        formatter: function(value, row, index) {
            var html = "<select class='form-table-product-template-attributte-values'>";
            if (row.AttributeValues && row.AttributeValues.length > 0) {

            }
            html += "</select>";
            return html;
        }
    }, {
        title: "操作",
        align: "center",
        field: 'action',
        formatter: function cellStyle(value, row, index) {
            var html = "";
            var url = "/product/attribute/line/";
            html += "<a class='form-tree-line-action-remove' href='#'><i class='fa fa-trash-o'></i></a>";
            return html;
        }
    }

], {
    onPostBody: function() {
        select2AjaxData(".form-table-product-template-attributte", '/product/attribute/?action=search');
        select2AjaxData(".form-table-product-template-attributte-values", "/product/attributevalue/?action=search", {
            queryParams: function(params) {
                var selectParams = {
                    name: params.term || "", // search term
                    offset: (params.page || 0) * LIMIT,
                    limit: LIMIT,
                };
                var xsrf = $("input[name ='_xsrf']");
                if (xsrf.length > 0) {
                    selectParams._xsrf = xsrf[0].value;
                }
                if ($(this).length > 0 && $(this)[0].nodeName == "SELECT") {
                    selectParams.exclude = $(this).val();
                }
                return selectParams
            }
        });
        formTreeLineActionRemove();

        $(".form-table-product-template-attributte .form-table-product-template-attributte-values").on("change", function(e) {
            console.log(e)
        });
    },
    queryParams: function(params) {
        var xsrf = $("input[name ='_xsrf']");
        if (xsrf.length > 0) {
            params._xsrf = xsrf[0].value;
        }
        params.action = 'table';

        var templateId = $("input[name ='recordID']");
        if (templateId.length > 0) {
            params.templateId = parseInt(templateId[0].value);
        } else {
            params.saleOrderId = 0;
        }
        return params;
    }
});
//产品款式的属性增加一行
$("#add-one-product-template-attribute").on("click", function(e) {
    $("#product-template-attribute-table").bootstrapTable("prepend", [{
        ID: null,
        Attribute: {
            id: null,
            Name: '',
        },
        AttributeValues: [{ id: null, Name: "" }]
    }]);
});