/*添加tab栏目*/
var tabChoseNum = [];
var tabNumz = 0;
function addPanel(titleName,href_,obj){
	if(href_ =="#"||href_ == "undefined")return;
	if(!$(obj).hasClass("active")){
		$(obj).addClass("active");
		tabChoseNum[tabNumz] = new Array()
		tabChoseNum[tabNumz].push(titleName);
		tabChoseNum[tabNumz].push($(obj).parent().attr("index"));
		tabNumz ++;
		$('#contentTab').tabs('add',{
			title: titleName,
			href: href_,
			closable: true
		});
	}else{
		$('#contentTab').tabs("select",parseInt($(obj).parent().attr("index"))+1)
	};
};
/*
function addPanel(titleName,href_,refreshTab){
	if(href_ =="#"||href_ == "undefined")return;
	if(!refreshTab){
		//$(obj).addClass("active");
		tabChoseNum[tabNumz] = new Array()
		tabChoseNum[tabNumz].push(titleName);
		tabChoseNum[tabNumz].push($(obj).parent().attr("index"));
		tabNumz ++;
		$('#contentTab').tabs('add',{
			title: titleName,
			href: href_,
			closable: true
		});
	}else{
		var tab = $('#contentTab').tabs('getSelected');
		$('#contentTab').tabs('update',{
			title: titleName,
			href: href_,
			closable: true
		});
		//$('#contentTab').tabs("select",parseInt($(obj).parent().attr("index"))+1)
	};
};
*/
/*设置tab内容区域高度*/
function setPanelHeight(){
	var wHeight = $(window).height(),
	headerHeight = $(".header").height(),
	tabHeight = $("#contentTab .tabs-header:first").height(),
	panelHeight = wHeight-headerHeight-tabHeight + 22;
	return panelHeight;
};
window.requestAnimFrame = (function(){  
        return  window.requestAnimationFrame       ||   
        window.webkitRequestAnimationFrame ||   
        window.mozRequestAnimationFrame    ||   
        window.oRequestAnimationFrame      ||   
        window.msRequestAnimationFrame     ||   
        function( callback ){  
            window.setTimeout(callback, 1000/60);  
        };  
})(); 

$.extend($.fn.validatebox.defaults.rules, {  
	ip: {// 验证IP地址
		validator: function (value) {
			var re=/^(\d+)\.(\d+)\.(\d+)\.(\d+)$/g //匹配IP地址的正则表达式
			if(re.test(value)){
				if( RegExp.$1 <256 && RegExp.$2<256 && RegExp.$3<256 && RegExp.$4<256) return true;
			}
			return false;
		},
		message: 'IP地址格式不正确'
	},       
	mac: {
		validator: function (value) {
			return isValidMacAddress(value);
		},
		message: 'mac地址不正确'
	}
}); 

/*mac地址验证*/
function isValidMacAddress(address) { 
   var macs = new Array();
    macs = address.split("-");
    if(macs.length != 6){
    	return false;
    }
	for(var s=0; s<6; s++) {
		var temp = parseInt(macs[s],16);
		if(isNaN(temp)){
			return false;
		};
		if(temp < 0 || temp > 255){
			return false;
		};
	};
    return true;
}
/*短信内容*/
function smsTxt(text){
	var d = new Date(),
		vYear = d.getFullYear(),
		vMon = d.getMonth() + 1,
		vDay = d.getDate();
	var today = vYear+"-"+vMon+"月"+vDay+"日";
	var text_ = text.replace(/\&用户名\&/g, "xxxx");
	text_ = text_.replace(/\&手机号\&/g, "00000000");
	text_ = text_.replace(/\&今日日期\&/g, today);
	return text_
}
$(function(){

	//var date = new Date();        
	//window.alert(new Date().pattern("yyyy-MM-dd"));   
})
   /**       
 * 对Date的扩展，将 Date 转化为指定格式的String       
 * 月(M)、日(d)、12小时(h)、24小时(H)、分(m)、秒(s)、周(E)、季度(q) 可以用 1-2 个占位符       
 * 年(y)可以用 1-4 个占位符，毫秒(S)只能用 1 个占位符(是 1-3 位的数字)       
 * eg:       
 * (new Date()).pattern("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423       
 * (new Date()).pattern("yyyy-MM-dd E HH:mm:ss") ==> 2009-03-10 二 20:09:04       
 * (new Date()).pattern("yyyy-MM-dd EE hh:mm:ss") ==> 2009-03-10 周二 08:09:04       
 * (new Date()).pattern("yyyy-MM-dd EEE hh:mm:ss") ==> 2009-03-10 星期二 08:09:04       
 * (new Date()).pattern("yyyy-M-d h:m:s.S") ==> 2006-7-2 8:9:4.18       
 */          
Date.prototype.pattern=function(fmt) {           
    var o = {           
    "M+" : this.getMonth()+1, //月份           
    "d+" : this.getDate(), //日           
    "h+" : this.getHours()%12 == 0 ? 12 : this.getHours()%12, //小时           
    "H+" : this.getHours(), //小时           
    "m+" : this.getMinutes(), //分           
    "s+" : this.getSeconds(), //秒           
    "q+" : Math.floor((this.getMonth()+3)/3), //季度           
    "S" : this.getMilliseconds() //毫秒           
    };           
    var week = {           
    "0" : "/u65e5",           
    "1" : "/u4e00",           
    "2" : "/u4e8c",           
    "3" : "/u4e09",           
    "4" : "/u56db",           
    "5" : "/u4e94",           
    "6" : "/u516d"          
    };           
    if(/(y+)/.test(fmt)){           
        fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length));           
    }           
    if(/(E+)/.test(fmt)){           
        fmt=fmt.replace(RegExp.$1, ((RegExp.$1.length>1) ? (RegExp.$1.length>2 ? "/u661f/u671f" : "/u5468") : "")+week[this.getDay()+""]);           
    }           
    for(var k in o){           
        if(new RegExp("("+ k +")").test(fmt)){           
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));           
        }           
    }           
    return fmt;           
}         
/*表格栏目展开\收缩*/      
function expandRowDisplay(obj){
	var $obj = obj;
	if($obj.hasClass('clicked')){
		/*展开*/
		$obj.removeClass('clicked');
		$obj.text('展开');
		$obj.parents('.datagrid-view').find('.datagrid-body>.datagrid-btable>tr[datagrid-row-index ="'+$obj.parent().parent().parent().attr('datagrid-row-index')+'"]').removeClass('expandRowActive');
	}else{		
		/*收缩*/
		$obj.addClass('clicked');
		$obj.text('收起');
		$obj.parents('.datagrid-view').find('.datagrid-body>.datagrid-btable>tr[datagrid-row-index ="'+$obj.parent().parent().parent().attr('datagrid-row-index')+'"]').addClass('expandRowActive');	
	};
	/*触发点击事件*/
	$obj.parents('.datagrid-view').find('.datagrid-body .datagrid-btable tr[datagrid-row-index ="'+$obj.parent().parent().parent().attr('datagrid-row-index')+'"] .datagrid-row-expander').trigger('click');
	
}

/*搜索栏目 下拉菜单操作列表*/
function searchMenu(){
	$('.searchBtnMenu .searchListTarget').bind("mouseenter",function(){
		$(this).parent().addClass("hover");		
	});
	$('.searchBtnMenu').bind("mouseleave",function(){
		$(this).removeClass("hover");		
	});
	$('.searchBtnMenu .opeBtnList li').bind("click",function(e){
		e.preventDefault();
		$(this).parent().parent().removeClass("hover");		
	});	
}

/*图片上传预览*/
jQuery.fn.extend({
	uploadPreview: function (opts) {
		var _self = this,
		_this = $(this);
		opts = jQuery.extend({
			Img: "ImgPr",
			Width: 100,
			Height: 100,
			ImgType: ["gif", "jpeg", "jpg", "bmp", "png"],
			imgMaxSize:1024000,
			Callback: function () {}
		}, opts || {});
		_self.getObjectURL = function (file) {
			var url = null;
			if (window.createObjectURL != undefined) {
				url = window.createObjectURL(file)
			} else if (window.URL != undefined) {
				url = window.URL.createObjectURL(file)
			} else if (window.webkitURL != undefined) {
				url = window.webkitURL.createObjectURL(file)
			}
			return url
		};
		_this.change(function () {
			if (this.value) {
				if (!RegExp("\.(" + opts.ImgType.join("|") + ")$", "i").test(this.value.toLowerCase())) {
					showMsg("选择文件错误,图片类型必须是" + opts.ImgType.join(",") + "中的一种",false, false, "warning");
					this.value = "";
					return false
				}
				_self.image = new Image();
     			_self.image.src = this.files[0];
     			var fileSize = 0;
				if($.browser.msie){
					fileSize = this.files[0].filesize;				
				}else{
					fileSize = this.files[0].size;				
				};
				if(fileSize > opts.imgMaxSize){
					showMsg("选择文件错误,图片大小必须小于"+ opts.imgMaxSize/1000+"k", false, false, "warning");
					return false
				};
				if ($.browser.msie) {
					try {
						$("#" + opts.Img).attr('src', _self.getObjectURL(this.files[0]))
					} catch (e) {
					var src = "";
					var obj = $("#" + opts.Img);
					var div = obj.parent("div")[0];
					_self.select();
					if (top != self) {
						window.parent.document.body.focus()
					} else {
						_self.blur()
					}
					src = document.selection.createRange().text;
					document.selection.empty();
					obj.hide();
					obj.parent("div").css({
						'filter': 'progid:DXImageTransform.Microsoft.AlphaImageLoader(sizingMethod=scale)',
						'width': opts.Width + 'px',
						'height': opts.Height + 'px'
					});
					div.filters.item("DXImageTransform.Microsoft.AlphaImageLoader").src = src
				}
				} else {
					$("#" + opts.Img).attr('src', _self.getObjectURL(this.files[0]))
				}
					opts.Callback(this.value)
			}
		})
	}
});

function opeMenuList(){
	$('.opeMenu').each(function(){
		var clone = $(this).children().clone();
		if($(this).children('a').length > 1){
			$(this).children('a:not(:first)').remove();
		};
	});
}