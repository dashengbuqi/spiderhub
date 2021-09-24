/**
 * 在文档加载过程中对金额进行处理
 * @param money
 */
function subMoney(money) {
	//alert("money:"+money);
	if(money == null || money == "" || money == undefined) {
		document.write("0.00");
	} else {
		document.write(money.substr(0, money.length-2));
	}
}

/**
 * 处理datagrid时，返回操作后的结果
 * @param money
 * @returns
 */
function subDatagridMoney(money) {
	if(money == '' || money == undefined || money == null) {
		return "0.00";
	} else {
		return money.substr(0, money.length-2);
	}
}

/**
 * 处理ajax返回的数据
 * 保留金额小数点后的两位
 */
function afterSubMoney(money) {
	//alert("money:"+money);
	var str = "";
	if(money == "" || money == undefined || money == null) {
		//alert(strMoney+",is null");
		str = "0.00";
	} else {
		var strMoney = money + "";
		//alert("strMoney="+strMoney);
		var len = strMoney.length - strMoney.indexOf(".") - 1;   // 用于记录有多少位的小数
		var index = strMoney.indexOf(".");  					// 记录小数点的位置
		//alert("len:"+len+",index:"+index);
		if(index < 0) {
			str = strMoney += ".00";
		} else if(len < 2) {
			str = strMoney += "0";
		} else if(len >= 2) {
			str = strMoney.substr(0, strMoney.length - len + 2);
		}
	}
	//alert("str is "+str);
	
	//alert("end");
	return str;
}