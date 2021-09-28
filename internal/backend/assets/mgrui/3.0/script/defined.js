$(function() {
	InitLeftMenu();
	tabClose();
	tabCloseEven();

});
// 初始化左侧

function InitLeftMenu() {

	$('.easyui-accordion li a').click(function() {
		$('.easyui-accordion li div').removeClass("selected");
		$(this).parent().addClass("selected");
	}).hover(function() {
		$(this).parent().addClass("hover");
	}, function() {
		$(this).parent().removeClass("hover");
	});

}

// 获取左侧导航的图标
function getIcon(menuid) {
	var icon = 'icon ';
	$.each(_menus.menus, function(i, n) {
		$.each(n.menus, function(j, o) {
			if (o.menuid == menuid) {
				icon += o.icon;
			}
		});
	});

	return icon;
}
/**
 * 刷新tab subtitle：标题 url：请求
 */
function refreshTab(subtitle, url) {
	var obj = getTabObj();
	var ifrmaeEl =obj.tabs('getSelected');
	var iframe = ifrmaeEl.find("iframe[name=" + subtitle + "]");
	destroyIframe(iframe);
	if (obj.tabs('exists', subtitle)) {
		var currTab = obj.tabs('getTab', subtitle);
		if (typeof (url) == 'undefined')
			url = $(currTab.panel('options').content).attr('src');
		obj.tabs('update', {
			tab : currTab,
			options : {
				content : createFrame(url, subtitle)
			}
		});
	}
}

function refreshCurrTab() {
	var obj = getTabObj();
	var currTab = obj.tabs('getSelected');
	var url = $(currTab.panel('options').content).attr('src');
	var subtitle = currTab.panel('options').title;
	obj.tabs('update', {
		tab : currTab,
		options : {
			content : createFrame(url, subtitle)
		}
	});
}
/** 获取标签 */
function getTabObj() {
	var obj = $('#tabs');
	if (obj.length == 0) {
		obj = top.parent.$('#tabs');
	}
	return obj;
}
function closeCurrTab() {
	var currTab = getTabObj().tabs('getSelected');
	var currtab_title = currTab.panel('options').title;
	obj.tabs('close', currtab_title);

}

/**
 * subtitle：tab名称 url：请求地址 ico：tab小图标 closable：是否允许关闭table
 * subName 父元素的name  更新tab用
 */
function addTab(subtitle, url, icon, closable,menuId,subName) {
	if (typeof (closable) == "undefined") {
		closable = true;
	}
	var tabs = getTabObj();
	// 如果tab不存在
	if (!tabs.tabs('exists', subtitle)) {
		tabs.tabs('add', {
			tabHeight : 40,
			title : subtitle,
			content : createFrame(url, subtitle,subName),
			closable : closable,
			cache : false,
			icon : icon
		});
	} else {
		// refreshCurrTab();
		tabs.tabs('select', subtitle);
		refreshTab(subtitle, url);
	}
	if (closable) {
		// 增加关闭事件
		tabClose();
	}
	if(menuId != undefined){
		$(".leftMenuTempClass").each(function(){
			$(this).removeClass("j_menu_active");
		});
		$("#leftMenu .menuLine a").each(function(){
			$(this).css("color","#363f45");  
		});
		$("#leftMenu"+menuId).addClass("j_menu_active");
		$("#leftMenu"+menuId).find("a").first().css("color","#055197");
		
		$("#tabs .tabs .tabs-selected").attr("menuid",menuId);
		
		$(".left_up1").each(function(){
			$(this).removeClass("left_up1");
			$(this).addClass("left_up3");
		});
		
		//获取选中菜单的父菜单ID
		var parentId =  $("#leftMenu"+menuId).attr("parentid");
		$("#leftMenu"+parentId).removeClass("left_up3");
		$("#leftMenu"+parentId).addClass("left_up1");
		
		$("#tabs .tabs li").unbind("click");  
		$("#tabs .tabs li").bind("click",function(){
			//取消左侧导航的选择样式
			$(".leftMenuTempClass").each(function(){
				$(this).removeClass("j_menu_active");
			});
			$("#leftMenu .menuLine a").each(function(){
				$(this).css("color","#363f45");  
			});
			
			//获取当前点击的menuId
			var menuId = $(this).attr("menuid");
			//给选中的菜单增加样式
			$("#leftMenu"+menuId).addClass("j_menu_active");
			$("#leftMenu"+menuId).find("a").first().css("color","#055197");
			
			$(".left_up1").each(function(){
				$(this).removeClass("left_up1");
				$(this).addClass("left_up3");
			});
			
			//获取选中菜单的父菜单ID
			var parentId =  $("#leftMenu"+menuId).attr("parentid");
			//判断父菜单是否隐藏
			if($("#end"+parentId).css("display")=='none') {
				$("#end"+parentId).slideToggle(200);
        		$("#up"+parentId).attr("class","left_up1");
        	}else{
        		$("#leftMenu"+parentId).removeClass("left_up3");
    			$("#leftMenu"+parentId).addClass("left_up1");
        	} 
			
			//获取顶层菜单id
			var topMenuId = $("#leftMenu"+parentId).attr("parentid");
			//判断当前选中的菜单父菜单是否显示
			if($("#leftMenu"+parentId).is(':hidden')){
				var topMenuName = $("#topMenu"+topMenuId).attr("menuname");
				$('#west').panel({title:topMenuName});
				selectMenu(topMenuId);
			}

			/**
			 * 选中一级菜单
			 */
			$(".topMenuTempClass").each(function(){
        		$(this).find("a").first().removeClass("j_menu_top_active");
        	});
			$("#topMenu"+topMenuId).find("a").first().addClass("j_menu_top_active");
		});
	}
}

function createFrame(url, subtitle,subName) {
	// 为url增加提交参数addTabFlag，说明当前请求是新打开的一个页签
	if (url.indexOf("?") >= 0) {
		url = url + "&addTabFlag=1";
	} else {
		url = url + "?addTabFlag=1";
	}
	url = url + "&t=" + Math.random();
	var s = '<iframe scrolling="auto" frameborder="0" parent_name="'+subName+'" name="' + subtitle
			+ '" src="' + url + '" style="width:100%;height:100%;"></iframe>';
	return s;
}
/* 日程tab */
function addTabForDate(subtitle, url, icon, closable) {
	if (typeof (closable) == "undefined") {
		closable = true;
	}
	var tabs = getTabObj();
	// 如果tab不存在
	if (!tabs.tabs('exists', subtitle)) {
		tabs.tabs('add', {
			tabHeight : 40,
			title : subtitle,
			content : createFrameForDate(url, subtitle),
			closable : closable,
			cache : false,
			icon : icon
		});
	} else {
		// refreshCurrTab();
		tabs.tabs('select', subtitle);
		refreshTab(subtitle, url);
	}
	if (closable) {
		// 增加关闭事件
		tabClose();
	}
}

function createFrameForDate(url, subtitle) {
	// 为url增加提交参数addTabFlag，说明当前请求是新打开的一个页签
	if (url.indexOf("?") >= 0) {
		url = url + "&addTabFlag=1";
	} else {
		url = url + "?addTabFlag=1";
	}
	url = url + "&t=" + Math.random();
	var s = '<iframe scrolling="no" frameborder="0" name="' + subtitle
			+ '" src="' + url + '" style="width:100%;height:100%;"></iframe>';
	return s;
}

function tabClose() {
	/* 双击关闭TAB选项卡 */
	$(".tabs-inner").dblclick(function() {
		var subtitle = $(this).children(".tabs-closable").text();
		closeTab(subtitle);
	});
	/* 为选项卡绑定右键 */
	$(".tabs-inner").bind('contextmenu', function(e) {
		$('#mm').menu('show', {
			left : e.pageX,
			top : e.pageY
		});
		var subtitle = $(this).children(".tabs-closable").text();
		$('#mm').data("currtab", subtitle);
		$('#tabs').tabs('select', subtitle);
		return false;
	});
}
/**
 * 关闭tab
 * 
 * @param subtitle
 *            tab的标题
 */
function closeTab(subtitle) {
	var ifrmaeEl = getTabObj().tabs('getSelected');
	var iframe = ifrmaeEl.find("iframe[name=" + subtitle + "]");
	destroyIframe(iframe);
	$('#tabs').tabs('close', subtitle);
	//$('#tabs').tabs('close', subtitle);
}

/*清除iframe防止长时间过长，页面变慢*/
function destroyIframe(iframe){
	var iframes = $(iframe).prop('contentWindow');

	$(iframe).attr('src', 'about:blank');
	try{
		iframes.document.write('');
		iframes.document.clear();
	}catch(e){}
	//把iframe从页面移除
	$(iframe).remove();

}
// 绑定右键菜单事件
function tabCloseEven() {
	// 刷新
	$('#mm-tabupdate').click(function() {
		var currtab_title = $('#mm').data("currtab");
		var ifrmaeEl = getTabObj().tabs('getSelected');
		var iframe = ifrmaeEl.find("iframe[name=" + currtab_title + "]");
		destroyIframe(iframe);
		var currTab = $('#tabs').tabs('getSelected');
		var url = $(currTab.panel('options').content).attr('src');
		$('#tabs').tabs('update', {
			tab : currTab,
			options : {
				content : createFrame(url)
			}
		});
	});
	// 关闭当前
	$('#mm-tabclose').click(function() {
		var currtab_title = $('#mm').data("currtab");
		closeTab(currtab_title);
	});
	// 全部关闭
	$('#mm-tabcloseall').click(function() {
		$('.tabs-inner span').each(function(i, n) {
			var t = $(n).text();
			if ($(n).attr("class").indexOf("tabs-closable") > 0)
		         	closeTab(t);
		});
	});
	// 关闭除当前之外的TAB
	$('#mm-tabcloseother').click(function() {
		$('#mm-tabcloseright').click();
		$('#mm-tabcloseleft').click();
	});
	// 关闭当前右侧的TAB
	$('#mm-tabcloseright').click(
			function() {
				var nextall = $('.tabs-selected').nextAll();
				if (nextall.length == 0) {
					return false;
				}
				nextall.each(function(i, n) {
					var t = $('a:eq(0) span', $(n)).text();
					if ($('a:eq(0) span', $(n)).attr("class").indexOf(
							"tabs-closable") > 0)
						closeTab(t);
				});
				return false;
			});
	// 关闭当前左侧的TAB
	$('#mm-tabcloseleft').click(
			function() {
				var prevall = $('.tabs-selected').prevAll();
				if (prevall.length == 0) {
					return false;
				}
				prevall.each(function(i, n) {
					var t = $('a:eq(0) span', $(n)).text();
					if ($('a:eq(0) span', $(n)).attr("class").indexOf(
							"tabs-closable") > 0)
						 closeTab(t);
				});
				return false;
			});

	// 退出
	$("#mm-exit").click(function() {
		$('#mm').menu('hide');
	});
}
// 格式化时间控件的日期格式：yyyy-mm-dd
$.fn.datebox.defaults.formatter = function(date) {
	var y = date.getFullYear();
	var m = date.getMonth() + 1;
	var d = date.getDate();
	return y + '-' + (m < 10 ? ('0' + m) : m) + '-' + (d < 10 ? ('0' + d) : d);
};
$.fn.datebox.defaults.parser = function(s) {
	if (!s)
		return new Date();
	var ss = s.split('-');
	var y = parseInt(ss[0], 10);
	var m = parseInt(ss[1], 10);
	var d = parseInt(ss[2], 10);
	if (!isNaN(y) && !isNaN(m) && !isNaN(d)) {
		return new Date(y, m - 1, d);
	} else {
		return new Date();
	}
};

/** *常量** */
var LOAD_MSG = "正在加载消息...";
var DATAGRID_ID = "#datagrid";

function getDatagridWidth() {
	return parseInt($(window).width()) - 1;
}

function getDatagridHeight() {
	var h = 0;
	if ($(".searArea").css("display") == 'block') {
		h += $(".searArea").outerHeight(true);
	}
	if ($(".shouwbtn").css("display") == 'block') {
		h += $(".shouwbtn").outerHeight(true);
	}
	if ($(".total").css("display") == 'block') {
		h += $(".total").outerHeight(true);
	}
	return parseInt($(window).height()) - h - 1;
}
var lastWidth = $(window).height();

$(document).ready(function() {
	$(window).resize(function() {
		if (lastWidth > $(window).height())
			resizeDataGrid(18, 0);
		else {
			resizeDataGrid(0, 18);
		}
	});
	addLegendClick();

	$('.search_close').click(function() {
		closeOrOpen();
	});
	//showUserSearchCondition();
	loadshortcutSearch();

});
/***/
function resizeDataGrid(w1, w2) {
	if ($(DATAGRID_ID)) {
		$(DATAGRID_ID).datagrid('resize', {
			width : getDatagridWidth() - w1,
			height : getDatagridHeight()
		});
		$(DATAGRID_ID).datagrid('resize', {
			width : getDatagridWidth() - w2,
			height : getDatagridHeight()
		});
	}
}

/**
 * @param {}
 *            dataUrl
 */
function GoEnterPage() {
	var e = jQuery.Event("keydown");
	e.keyCode = 13;
	$("input.pagination-num").trigger(e);
}
/**
 * 重新加载datagrid
 * 
 * @param {}
 *            dataUrl eg: /user/list.html?data=1&userAccount=jim&userName=陈
 */
function reloadDatagrid(dataUrl) {
	if (typeof (dataUrl) != 'undefined') {
		var urlTmp = encodeURI(dataUrl);
		//%252F  => /
		urlTmp = urlTmp.replace(/%252F/g, '/');
		$(DATAGRID_ID).datagrid('options').url = urlTmp;
	}
	$(DATAGRID_ID).datagrid('reload');
}

function reloadTreegrid(dataUrl) {
	if (typeof (dataUrl) != 'undefined') {
		$(DATAGRID_ID).treegrid('options').url = encodeURI(dataUrl);
	}
	$(DATAGRID_ID).treegrid('reload');
}
/**
 * 刷新dialog_common.jsp
 * 
 * @param {}
 *            dataUrl
 */
function reloadDialogCommon(dataUrl) {
	if (typeof (dataUrl) != 'undefined') {
		$("#dialogCommonId").panel({
			region : 'center',
			href : encodeURI(dataUrl)
		});
	}
	// $("#dialogCommonId").panel('refresh');
}

/**
 * 给legend添加事件
 */
function addLegendClick() {
	$("legend").each(function() {
		$(this).click(function() {
			if ($("span", $(this)).attr("class") == "presentation") {
				$("span", $(this)).attr("class", "presentation_down");
				$($(this)).next().slideToggle(100);
				$(this).addClass("sp_heigth");
				$(this).parent().attr("style", "padding-bottom:0px");
			} else {
				$("span", $(this)).attr("class", "presentation");
				$(this).removeClass("sp_heigth");
				$($(this)).next().toggle(100);
				$(this).parent().attr("style", "padding-bottom:10px");
			}
		});
	});
}
/**
 * 日期格式化
 * 
 * @param {}
 *            val
 * @param isDate
 *            是否日期，日期没有时分秒
 * @return {}
 */
function formattime(val, isDate) {
	if (val != null && val != '') {
		var year = parseInt(val.year) + 1900;
		var month = (parseInt(val.month) + 1);
		month = month > 9 ? month : ('0' + month);
		var date = parseInt(val.date);
		date = date > 9 ? date : ('0' + date);
		if (typeof (isDate) == 'undefined' || !isDate) {
			var hours = parseInt(val.hours);
			hours = hours > 9 ? hours : ('0' + hours);
			var minutes = parseInt(val.minutes);
			minutes = minutes > 9 ? minutes : ('0' + minutes);
			var seconds = parseInt(val.seconds);
			seconds = seconds > 9 ? seconds : ('0' + seconds);
			return year + '/' + month + '/' + date + ' ' + hours + ':'
					+ minutes + ':' + seconds;
		}
		return year + '/' + month + '/' + date;
	}
	return '';
}

/**
 * 日期格式化
 * 
 * @param {}
 *            val
 * @param isDate
 *            是否日期，日期有时分秒
 * @return {}
 */
function formatDate(val, isDate) {
	if (val != null && val != '') {
		var year = parseInt(val.year) + 1900;
		var month = (parseInt(val.month) + 1);
		month = month > 9 ? month : ('0' + month);
		var date = parseInt(val.date);
		date = date > 9 ? date : ('0' + date);
		if (typeof (isDate) == 'undefined' || !isDate) {
			var hours = parseInt(val.hours);
			hours = hours > 9 ? hours : ('0' + hours);
			var minutes = parseInt(val.minutes);
			minutes = minutes > 9 ? minutes : ('0' + minutes);
			var seconds = parseInt(val.seconds);
			seconds = seconds > 9 ? seconds : ('0' + seconds);
			return year + '-' + month + '-' + date + ' ' + hours + ':'
			+ minutes + ':' + seconds;
		}
		return year + '-' + month + '-' + date;
	}
	return '';
}

/**
 * 进度条 showProcess(true, '温馨提示', '正在提交数据...'); showProcess(false); 关闭进度条
 * 
 */
function showProcess(isShow, title, msg) {
	if (!isShow) {
		$.messager.progress('close');
		return;
	}
	$.messager.progress({
		title : title,
		msg : msg
	});
}
/**
 * 信息提示 refreshDatagrid：传true过来则刷新Datagrid refreshTab：传true过来则刷新当前tab
 * icon四种设置："error"、"info"、"question"、"warning"
 */
function showMsg(msg, refreshDatagrid, refreshTab, icon) {
	if (typeof (icon) == 'undefined')
		icon = "info";
	$.messager.alert('操作提示', msg, icon, function() {
		if (refreshDatagrid) {
			try {
				parent.reloadDatagrid();
				reloadDatagrid();
				parent.closeMyWindow();
			} catch (e) {
			}
		}
		if (refreshTab) {
			refreshCurrTab();
		}

	});
}
/**
 * 确认提示框
 * 
 * @param {}
 *            title
 * @param {}
 *            msg
 * @param {}
 *            callback 确认后执行的方法
 */
function showConfirm(title, msg, callback) {
    
	var messagebox= $.messager.confirm(title, msg, function(r) {
		if (r) {
			if (jQuery.isFunction(callback))
				callback.call();
		}
	});
        return messagebox;
}

$(function() {
	$('body').append(
			'<div id="myWindow" class="easyui-dialog" closed="true"></div>');
});

$.fn.my_openbox=function(title,url,width,height){

	width = width === undefined ? 550 : width;
	height = height === undefined ? 500 : height;

	if ($('#myWindow').length == 0) {
		$('body')
			.append(
				'<div id="myWindow" class="easyui-dialog" closed="true"></div>');
	}
	$('#myWindow').dialog({
		title: title,
		width: width,
		height: height,
		closed: false,
		cache: false,
		href: url,
		onOpen:function(){
			var closeBox =$('.window').find(".panel-tool-close");
			closeBox.unbind("click").bind("click",function(){
				closeMyWindow();
			});
		},
		onClose : function() {
			$(this).dialog('destroy');
		},
		modal: true
	});
};


/**
 * 打开窗口
 * 
 * @param {}
 *            title
 * @param {}
 *            href
 * @param {}
 *            isReadonly 等于true时只读
 * @param {}
 *            width
 * @param {}
 *            height
 * @param {}
 *            modal
 * @param {}
 *            minimizable
 * @param {}
 *            maximizable
 */
function showMyWindow(title, href, isReadonly, width, height, modal,
		minimizable, maximizable) {
	if ($('#myWindow').length == 0) {
		$('body')
				.append(
						'<div id="myWindow" class="easyui-dialog" closed="true"></div>');
	}
	var parm = "";
	if (typeof (isReadonly) != 'undefined' && isReadonly == true) {
		parm = "&isReadonly=1";
	}
	var bodyWidth = $("body").width() * 0.8;
	var bodyHeight = $("body").height() * 0.9;

	$('#myWindow')
			.window(
					{
						title : title,
                                                onOpen:function(){
                                                    var closeBox =$('.window').find(".panel-tool-close");
                                                    closeBox.unbind("click").bind("click",function(){
                                                        closeMyWindow();
                                                    });                                                    
                                                },
						width : width === undefined ? bodyWidth : width,
						height : height === undefined ? bodyHeight : height,
						content : '<iframe id="showMyWindowId" scrolling="yes" frameborder="0"  src=/default/dialog?1=1'
								+ parm
								+ '&url='
								+ encodeURIComponent(href)
								+ ' style="width:100%;height:98%;"></iframe>',
						modal : modal === undefined ? true : modal,
						minimizable : minimizable === undefined ? false
								: minimizable,
						maximizable : maximizable === undefined ? false
								: maximizable,
						shadow : false,
						cache : false,
						closed : false,
						collapsible : false,
						resizable : false,
						loadingMessage : '正在加载数据，请稍等片刻......'
					});
}
/**
 * 打开审核窗口
 * 
 * @param {}
 *            title
 * @param {}
 *            href
 * @param {}
 *            isReadonly 等于true时只读
 * @param {}
 *            width
 * @param {}
 *            height
 * @param {}
 *            modal
 * @param {}
 *            minimizable
 * @param {}
 *            maximizable
 */
function showCheckWindow(title, href, isReadonly, width, height, institutionId,
		modal, minimizable, maximizable) {
	if ($('#myWindow').length == 0) {
		$('body')
				.append(
						'<div id="myWindow" class="easyui-dialog" closed="true"></div>');
	}
	var parm = "";
	if (typeof (isReadonly) == 'undefined' || isReadonly == true) {
		parm = "&isReadonly=1";
	}
	var institutionIdParam = "&institutionId=" + institutionId;
	$('#myWindow')
			.window(
					{
						title : title,
						width : width === undefined ? 600 : width,
						height : height === undefined ? 500 : height,
						content : '<iframe id="showMyWindowId" scrolling="yes" frameborder="0"  src="/systemUI/dialog_check.jsp?1=1'
								+ institutionIdParam
								+ ''
								+ parm
								+ '&url='
								+ encodeURIComponent(href)
								+ '" style="width:100%;height:98%;"></iframe>',
						modal : modal === undefined ? true : modal,
						minimizable : minimizable === undefined ? false
								: minimizable,
						maximizable : maximizable === undefined ? false
								: maximizable,
						shadow : false,
						cache : false,
						closed : false,
						collapsible : false,
						resizable : false,
						loadingMessage : '正在加载数据，请稍等片刻......'
					});
}
function closeMyWindow() {
	var iframe = $('#myWindow').find("iframe");
	destroyIframe(iframe);
	$('#myWindow').window('destroy');

}

/**
 * 打开分配管理员窗口
 * 
 * @param {}
 *            title
 * @param {}
 *            href
 * @param {}
 *            isReadonly 等于true时只读
 * @param {}
 *            width
 * @param {}
 *            height
 * @param {}
 *            modal
 * @param {}
 *            minimizable
 * @param {}
 *            maximizable
 */
function showAdminWindow(title, href, isReadonly, width, height, modal,
		minimizable, maximizable) {
	if ($('#myWindow').length == 0) {
		$('body')
				.append(
						'<div id="myWindow" class="easyui-dialog" closed="true"></div>');
	}
	var parm = "";
	if (typeof (isReadonly) == 'undefined' || isReadonly == true) {
		parm = "&isReadonly=1";
	}
	var bodyWidth = $("body").width() * 0.8;
	var bodyHeight = $("body").height() * 0.9;
	$('#myWindow')
			.window(
					{
						title : title,
						width : width === undefined ? bodyWidth : width,
						height : height === undefined ? bodyHeight : height,
						content : '<iframe id="showMyWindowId" scrolling="yes" frameborder="0"  src="/systemUI/dialog_admin.jsp?1=1'
								+ parm
								+ '&url='
								+ encodeURIComponent(href)
								+ '" style="width:100%;height:98%;"></iframe>',
						modal : modal === undefined ? true : modal,
						minimizable : minimizable === undefined ? false
								: minimizable,
						maximizable : maximizable === undefined ? false
								: maximizable,
						shadow : false,
						cache : false,
						closed : false,
						collapsible : false,
						resizable : true,
						loadingMessage : '正在加载数据，请稍等片刻......'
					});
}
function closeAdminWindow() {
	$('#myWindow').window('destroy');
}
/**
 * 获得选中行的ID
 * 
 * @return {}
 */
function getSelected(func) {
	var selected = $(DATAGRID_ID).datagrid('getSelections');
	var ids = "";
	for ( var i = 0; i < selected.length; i++) {
		if (ids != "")
			ids += ",";
		// ids += selected[i].userId;
		ids += func(selected[i]);
	}
	return ids;
}
$(function() {
	//initDatagrid();
	/* 搜索栏目 下拉菜单操作列表 */
	$('.searchBtnMenu .searchListTarget').bind("mouseenter", function() {
		$(this).parent().addClass("hover");
	});
	$('.searchBtnMenu').bind("mouseleave", function() {
		$(this).removeClass("hover");
	});
	$('.searchBtnMenu .opeBtnList li').bind("click", function() {
		$(this).parent().parent().parent().removeClass("hover");
	});
});

function dataGridOnLoadSuccess(data){
		df_pagination($(DATAGRID_ID).datagrid('getPager'));
		addMenuForSearch();
		$(DATAGRID_ID).datagrid('unselectAll');
		$(DATAGRID_ID).datagrid('collapseGroup');
    	$(DATAGRID_ID).datagrid('cancelCellTip');
		$(DATAGRID_ID).datagrid('doCellTip', {
			onlyShowInterrupt : true,
			position : 'bottom',
			tipStyler : {
				'border' : '1px solid #333',
				'padding' : '2px',
				'color' : '#333',
				'background' : '#f7f5d1',
				'position' : 'absolute',
				'max-width' : '100%',
				'border-radius' : '4px',
				'-moz-border-radius' : '4px',
				'-webkit-border-radius' : '4px',
				boxShadow : '1px 1px 3px #292929'
			},
			delay : 100
		});	
}


function initDatagrid() {

	if ($(DATAGRID_ID).length > 0) {
		$(DATAGRID_ID).datagrid({
			width : getDatagridWidth(),
			height : getDatagridHeight(),
			url : dataUrl,
			fitColumns : true,
			nowrap : true,
			autoRowHeight : false,
			striped : true,
			collapsible : true,
			type : "post",
			pagination : true,
			pageSize : 15,// 每页记录数
			pageList : [ 15, 30, 45 ], // 分页记录数数组
			rownumbers : true,
			checkOnSelect : false,
			selectOnCheck : false,
			singleSelect : true,
			/* onSortColumn: function (sort, order) {}, */
			// datagrid加载完后，调用以下方法
			onLoadSuccess : function(data) {
				df_pagination($(DATAGRID_ID).datagrid('getPager'));
				addMenuForSearch();
				$(DATAGRID_ID).datagrid('unselectAll');
				$(DATAGRID_ID).datagrid('doCellTip', {
					onlyShowInterrupt : true,
					position : 'bottom',
					tipStyler : {
						'border' : '1px solid #333',
						'padding' : '2px',
						'color' : '#333',
						'background' : '#f7f5d1',
						'position' : 'absolute',
						'max-width' : '100%',
						'border-radius' : '4px',
						'-moz-border-radius' : '4px',
						'-webkit-border-radius' : '4px',
						boxShadow : '1px 1px 3px #292929'
					},
					delay : 100
				});
			}
		});
		fieldsetClick();
	}
}

function df_pagination(obj) {
	obj
			.pagination({
				showRefresh : true,
				beforePageText : "第",
				afterPageText : "页 <a href='javascript:void(0)' onclick='GoEnterPage()'>GO</a>，共{pages}页",
				displayMsg : '当前 {from} 到 {to} 条，总共 {total} 条&nbsp;&nbsp;&nbsp;&nbsp;'
			});
	
}

function selectMoreColumn() {
	var obj = $("#tzl");
	var offset = obj.offset();
	var comBox = $(".comBox");
	if (comBox.css("display") == 'block') {
		comBox.css("display", "none");
	} else {
		if ($('.comBox').html().length < 100) {
			$.ajax({
				type : 'POST',
				async : true,
				url : "/column/showDefault.html",
				dataType : "html",
				error : function() {
					alert("系统异常！");
				},
				success : function(html) {
					$('.comBox').html(html);
				}
			});
		}
		comBox.css("left", offset.left + obj.width() - comBox.width());
		comBox.css("top", offset.top + obj.height() - comBox.height() - 10);
		comBox.css("display", "block");
	}
}

/**
 * 提交form表单
 * 
 * @param {}
 *            dataUrl
 * @param {}
 *            func
 */
function submitForm(dataUrl, func, formId, data) {
	if (typeof (formId) == 'undefined')
		formId = 'datagridForm';
	$('#' + formId).form('submit', {
		url : dataUrl,
		onSubmit : function() {
			var flag = $(this).form('validate');
			if (flag) {
				showProcess(true, '温馨提示', '正在提交数据...');
			}
			return flag;
		},
		success : function(data) {
			showProcess(false);
			func(data);
		},
		onLoadError : function() {
			showProcess(false);
			$.messager.alert("温馨提示', '由于网络或服务器太忙，提交失败，请重试！");
		}
	});
}
var tzFlag = false; // 拖拽标识
var field = [];
/** 获取datagrid选中的checkbox对应行指定field的值 */
function getChecked(field) {
	var checkedItems = $(DATAGRID_ID).datagrid('getChecked');
	var arr = [];
	$.each(checkedItems, function(index, item) {
		arr.push(eval("item." + field));
	});
	return arr.join(",");
}
/**
 * 给列表增加搜索条件
 */
function addMenuForSearch() {

	$(".datagrid-header-row > td").each(function() {
		if ($(this).attr("field") != null && $(this).attr("field") != 'opt') {
			$(this).bind('contextmenu', function(e) {
				field[0] = $(this).attr("field");
				field[1] = $(this).find('span').html();
				$("#searchMenu").menu('show', {
					left : e.pageX,
					top : e.pageY
				});
				return false;
			});
		}
	});

	$("#s_condition").click(function() {
		addSearchCondition();
		resizeDataGrid(18, 0);
	});
}

document
		.writeln("<div id=\"searchMenu\" class=\"easyui-menu\" style=\"display:none;width:150px;\">"
				+ "<div id=\"s_condition\" data-options=\"iconCls:'icon-search'\">添加查询条件</div>"
				+ "</div>");

/**
 * 增加查询条件
 */
function addSearchCondition() {
	if (field[0] == 'opt')
		return;
	var flag = false;
	$(".searArea ul li div").each(
			function() {

				var temp = field[0].substring(0, 1).toUpperCase()
						+ field[0].substring(1);
				var searchConditionName = $(this).find("input,select").attr(
						"name");
				if (field[0] == searchConditionName
						|| ("start" + temp) == searchConditionName) {
					flag = true;
				}
			});
	if (!flag) {
		if (costomSearchCondition(field)) {
			addConditionHtml(field);
			$.ajax({
				type : 'POST',
				async : true,
				url : "/condition/user/add.html",
				data : "field=" + field[0] + "&label=" + field[1],
				dataType : "json"
			});
		}
	}
}
/** 通过重写次方法，自定义查询条件 */
function costomSearchCondition(arg) {
	return true;
}
/** 增加查询条件 */
function addConditionHtml(field) {
	field = eval(field);
	if (field[1].indexOf('时间') != -1) {
		var temp = field[0].substring(0, 1).toUpperCase()
				+ field[0].substring(1);
		$(".searArea ul li")
				.prepend(
						"<div class=\"s_c\" ><em class=\"deleteb\" onclick=\"deleteCondition(this,'"
								+ field[0]
								+ "')\"></em><label>"
								+ field[1]
								+ "：</label>"
								+ "<input id=\"start"
								+ temp
								+ "\" name=\"start"
								+ temp
								+ "\"  type=\"text\" placeholder=\"选择开始时间\" class=\"calBg wid140\" onfocus=\"WdatePicker({startDate:'%y-%M-%d 00:00:00',dateFmt:'yyyy-MM-dd HH:mm:ss'})\" />  -  "
								+ "<input id=\"end"
								+ temp
								+ "\" name=\"end"
								+ temp
								+ "\" type=\"text\" placeholder=\"选择结束时间\" class=\"calBg wid140\" onfocus=\"WdatePicker({startDate:'%y-%M-%d 00:00:00',dateFmt:'yyyy-MM-dd HH:mm:ss'})\"/></div>");

	} else {
		$(".searArea ul li").prepend(
				"<div class=\"s_c\" ><em class=\"deleteb\" onclick=\"deleteCondition(this,'"
						+ field[0] + "')\"></em><label>" + field[1]
						+ "：</label>  <input name=\"" + field[0]
						+ "\" type=\"text\" class=\"wid140\" /></div>");
	}
}

function deleteCondition(obj, field) {
	obj = $(obj);
	obj.parent().remove();
	$.ajax({
		type : 'POST',
		async : true,
		url : "/condition/user/delete.html",
		data : "field=" + field,
		dataType : "json"
	});
}

/**
 * 显示用户自定义搜索条件 needShowUserCondition 显示用户搜索条件的标识符，true为显示
 */
/*function showUserSearchCondition() {
	if (typeof (needShowUserCondition) != "undefined") {
		$.ajax({
			type : 'POST',
			async : false,
			url : "/condition/user/select.html",
			data : "1=1",
			dataType : "json",
			error : function() {
				alert("系统异常！");
			},
			success : function(json) {
				var arr = eval(json);
				for ( var i = 0; i < arr.length; i++) {
					var field = [];
					field[0] = arr[i].field;
					field[1] = arr[i].label;
					addConditionHtml(field);
				}
				resizeDataGrid(18, 0);
			}
		});
	}
}*/

/** 搜索查询条件展开或者关闭 */
function closeOrOpen() {
	if ($(".searArea").css("display") == 'block') {
		$(".searArea").hide();
		$(".search_close").attr("class", "search_open");
		resizeDataGrid(0, 0);
	} else {
		$(".searArea").slideToggle(100);
		$(".search_open").attr("class", "search_close");
		setTimeout(function() {
			resizeDataGrid(18, 0);
		}, 201);
	}
	// resizeDataGrid(0,18);
}
/** 提交查询表单 onsubmit="return toSearch(this)" */
function toSearch(formObj) {
	var param = "", specharsFlag = false;
	$(":text,:checked", $(formObj)).each(function() {
		var v = $(this).val();
		if (notEmpty(v) && notEmpty($(this).attr("name"))) {
				param += "&" + $(this).attr("name") + "=" + v;
		}
	});
	$(":selected", $(formObj)).each(function() {
		if (notEmpty($(this).val())) {
			param += "&" + $(this).parent().attr("name") + "=" + $(this).val();
		}
	});
	$(":hidden", $(formObj)).each(function() {
		if (notEmpty($(this).val())) {
			param += "&" + $(this).attr("name") + "=" + $(this).val();
		}
	});
	if (!specharsFlag) {
		reloadDatagrid(formObj.action + param);
	}
	return false;
}


function toSearchTreeGrid(formObj) {
	var param = "", specharsFlag = false;
	$(":text,:checked", $(formObj)).each(function() {
		var v = $(this).val();
		if (notEmpty(v) && notEmpty($(this).attr("name"))) {
				param += "&" + $(this).attr("name") + "=" + v;
		}
	});
	$(":selected", $(formObj)).each(function() {
		if (notEmpty($(this).val())) {
			param += "&" + $(this).parent().attr("name") + "=" + $(this).val();
		}
	});
	$(":hidden", $(formObj)).each(function() {
		if (notEmpty($(this).val())) {
			param += "&" + $(this).attr("name") + "=" + $(this).val();
		}
	});
	if (!specharsFlag) {
		reloadTreegrid(formObj.action + param);
	}
	return false;
}


function notEmpty(value) {
	return (value != '' && typeof (value) != 'undefined') ? true : false;
}
/** 保存搜索条件 */
function saveSearchCondition(formName) {
	var formObj = $("form[name=" + formName + "]");
	var specharsFlag = false;
	var arr = [];
	$(":text,:checked", $(formObj)).each(
			function() {
				var v = $(this).val();
				if (notEmpty(v) && notEmpty($(this).attr("name"))) {
					if (checkSpechars(v)) {
						showMsg("存在特殊字符，添加失败！");
						specharsFlag = true;
					} else {
						var label = $(this).prev('label').html();
						if (typeof (label) == 'undefined')
							label = '';
						arr.push("{type:'2', text:'', value:'" + v
								+ "', name:'" + $(this).attr("name")
								+ "',label:'" + label + "'}");
					}
				}
			});
	$(":selected", $(formObj)).each(
			function() {
				if (notEmpty($(this).val())) {
					arr.push("{type:'1', text:'" + $(this).text()
							+ "', value:'" + $(this).val() + "', name:'"
							+ $(this).parent().attr("name") + "',label:'"
							+ $(this).parent().prev('label').html() + "'}");
				}
			});
	if (!specharsFlag) {
		if (arr.length > 0) {
			$.ajax({
				type : 'POST',
				async : true,
				url : "/shortcut/save.html",
				data : "jsonCondition=[" + arr.join(",") + "]",
				dataType : "html",
				success : function(data) {
					$("#shortcutSpan").html(data);
					$("#shortcutSpan").css("display", "");
					resizeDataGrid(18, 0);
				}
			});
		} else {
			showMsg("要保存的查询条件不能为空！");
		}
	}
}
function priceCalInfo(value){
	$("#calPriceText").val(value);
}
 
function checkSpechars(value) {
	var pattern = new RegExp("[~'!@#$%^&*()-+=]");
	if (value != "" && value != null) {
		if (pattern.test(value)) {
			return true;
		}
	}
	return false;
}

function loadshortcutSearch() {
	if ($("#shortcutSpan").length > 0) {
		$.ajax({
			type : 'POST',
			async : false,
			url : "/shortcut/select.html",
			data : "",
			dataType : "html",
			success : function(data) {
				$("#shortcutSpan").html(data);
				$("#shortcutSpan").css("display", "");
				resizeDataGrid(18, 0);
			}
		});
	}
}
// wwj添加以下内容
// ajax 提交请求
function delSubmit(dataUrl, func) {
	$.ajax({
		url : dataUrl,
		type : 'post',
		success : function(data) {
			func(data);
		},
		error : function() {
			$.messager.alert("温馨提示', '由于网络或服务器太忙，提交失败，请重试！");
		}
	});
}
// 显示消息窗口 并刷新指定Datagrid
function showMyMsg(msg, refreshDatagrid, closeParent, refreshTab, icon) {
	if (typeof (icon) == 'undefined')
		icon = "info";
	$.messager.alert('操作提示', msg, icon, function() {
		if (closeParent) {
			try {
				parent.reloadByDatagrid(refreshDatagrid);
				reloadByDatagrid(refreshDatagrid);
				parent.closeMyWindow();
			} catch (e) {
			}
		} else {
			reloadByDatagrid(refreshDatagrid);
		}
		if (refreshTab) {
			refreshCurrTab();
		}
	});
}
// 显示消息窗口 并关闭当前tab
function showMyMsgByTab(msg, closeParent, closeTab, tabTitle, iframeName, icon) {
	if (typeof (icon) == 'undefined')
		icon = "info";
	$.messager.alert('操作提示', msg, icon, function() {
		if (closeParent) {
			try {
				reloadTabByDatagrid(iframeName);
				closeByTabs(closeTab, tabTitle);
			} catch (e) {
			}
		}
	});
}
// 刷新指定datagrid
function reloadByDatagrid(datagrid) {
	$("#" + datagrid).datagrid('reload');
}
/* 关闭tabs */
function closeByTabs(tab, tabTitle) {
//	window.parent.$('#' + tab).tabs('close', tabTitle);
}
/* 关闭tab并刷新表格 */
function reloadTabByDatagrid(iframeName) {
	window.parent.$("iframe[name=" + iframeName + "]").contents().find(
			".pagination-load").trigger("click");
}
// wwj 以上内容
document
		.writeln('<div class="comBox clb"><div class="panel-loading">Loading...</div></div>');

/* 点击标题收缩、展开栏目内容 */
function fieldsetClick() {
	$(document).bind("click", function(e) {
		if ($(e.target).closest('legend').length != 0) {
			$(e.target).closest('legend').parent().toggleClass("hideBox");
		}
	});
}

/* 短信内容 */
function smsTxt(text) {
	var d = new Date(), vYear = d.getFullYear(), vMon = d.getMonth() + 1, vDay = d
			.getDate();
	var today = vYear + "-" + vMon + "月" + vDay + "日";
	var text_ = text.replace(/\&用户名\&/g, "xxxx");
	text_ = text_.replace(/\&手机号\&/g, "00000000");
	text_ = text_.replace(/\&今日日期\&/g, today);
	return text_;
}

/*
 * 功能:某行不能被选中,全选也不可用 obj:表对象, rowNum:行数
 */
function freezeRow(obj, rowNum) {
	$(obj).prev().prev().find('.datagrid-header .datagrid-htable').find(
			'input[type = checkbox]').attr('disabled', true);
	$(obj).prev().prev().find('.datagrid-body .datagrid-btable tr').eq(rowNum)
			.find('input[type = checkbox]').attr('disabled', true);
}

function bindItemTitleHide() {
	$('.itemTitle').bind("click", function(e) {
		$(this).parent().toggleClass('hideCon');
		$(this).next().toggleClass("hideBox");
	});
}

var isChrome = navigator.userAgent.toLowerCase().match(/chrome/) != null;
if (isChrome) {
	setTimeout(function(){$('input[type=text]').attr("autocomplete", "off");},1)
}
