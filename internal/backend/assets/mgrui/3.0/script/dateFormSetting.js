// JavaScript Document
function dateFormSetting(option){
	var self = this;
	option = $.extend({
		container:null,
		control:null,
		controlShow:true,
		width:1031,
		maxYear:new Date().getFullYear()+1,
		minYear:new Date().getFullYear(),
		colors:[{code:"a",color:"#08be4d"},{code:"b",color:"#ff8400"},{code:"c",color:"#f86d6d"},{code:"d",color:"#6761db"},{code:"e",color:"#edc31d"},{code:"f",color:"#db61af"},{code:"g",color:"#ab48df"},{code:"h",color:"#24977b"},{code:"i",color:"#43bdcf"},{code:"j",color:"#fd2a00"}],
		//settingOption:[{name:"淡季",code:"a"},{name:"旺季",code:"b"}],
		days: [['Sunday','日'], ['Monday','一'], ['Tuesday','二'], ['Wednesday','三'], ['Thursday','四'], ['Friday','五'], ['Saturday','六']],
		months: [['January','1月'], ['February','2月'], ['March','3月'], ['April','4月'], ['May','5月'], ['June','6月'], ['July','7月'], ['August','8月'], ['September','9月'], ['October','10月'], ['November','11月'], ['December','12月']],
		data : null,
		markData:null,
		colorSection:["#fff3f2","#f5ffee","#e9feff","#fcf2ff","#e7f5ff"],
		editCallback:function(){}
	},option||{});
	/*初始化参数*/
	self.init = function(){
		self.container = option.container,
		self.control = option.control,
		self.controlShow = option.controlShow,
		self.colors = option.colors,
		self.width = option.width,
		self.days = option.days,
		self.months = option.months,
		self.markData = option.markData,
		self.editCallback = option.editCallback,
		self.maxYear = option.maxYear,
		self.colorSection = option.colorSection,
		self.minYear = option.minYear;
		if(option.data){
			self.dataOption = option.data;
		}else{
			self.dataOption = option.colors;
		};
	};
	self.init();
 	var d_date = new Date();
	var _date = {
		year: d_date.getFullYear(),
		month: d_date.getMonth() + 1,
		day: d_date.getDate(),
		week: d_date.getDay()
	};
	self.minYear = d_date.getFullYear();
	/*取色器
		code:颜色代号
	*/
	self.getColor = function(code){
		for(var g = 0; g < self.colors.length;g++){
			if(self.colors[g].code == code){
				return g;
				break;
			};
		};
	};
	/*标记编号器
		code:颜色代号
	*/
	self.getCodeNum = function(code){
		for(var g = 0; g < self.dataOption.length;g++){
			if(self.dataOption[g].code == code){
				return g;
				break;
			};
		};
	};
	/*创建表格*/
	self.createdGrid = function(){
		var dayN = 0,section = 0,sNum = 0;
		for(var i = 0;i < 38; i++){
			var tr = $('<tr></tr>');
			tr.appendTo(self.calendarGrid);
			for(var t = 0; t < self.months.length + 1; t ++){
				var td = null,dateTxt = null;
				if(i == 0){
					if(t == 0){
						td = $('<th class="dayTitle"></th>');
						dateTxt = '星期';
					}else{
						td = $('<th class="monthTh"></th>');
						dateTxt = self.months[t-1][1];
					};
					
					td.text(dateTxt).appendTo(tr);
				}else{
					if(t == 0){
						td = $('<th class="dayTh"></th>');
						dateTxt = self.days[dayN][1];
						td.css("background",self.colorSection[sNum]);
						section += 1;
						if(section > 6){
							sNum>self.colorSection.length-2?sNum = 0:sNum++;
							section = 0;
						};
					}else{
						td = $('<td></td>');
						dateTxt = "";
					};
					td.text(dateTxt).appendTo(tr);
				};	
			};
			if(i > 0){
				dayN > self.days.length-2?dayN = 0:dayN ++;
			};
		};
	};
	/*搜索当前位置*/
	self.yearNum = function(year){
		var z = -1;
		for(var n = 0; n < self.dataBox.find("input").length; n ++){
			if( self.dataBox.find("input").eq(n).attr("year") == year){
				z = n;
				break;
			};
		};
		return z;
	}
	/*创建数据容器*/
	self.creatDataBox = function(year,value){
		var box = null;
		if(self.yearNum(year) == "-1"){
			 box = $('<input type="hidden" year="'+year+'" value="'+value+'" change="false"/>');
		};
		return box;
	};
	/*创建主体*/
	self.createdCalendar = function(){
		self.yy = _date.year;
		self.container.css({
			width:self.width
		}).empty();
		self.calendarYearMain = $('<div class="calendarYearHeader"></div>');
		self.prevBtn = $('<div class="calendarPrev"></div>');
		self.nextBtn = $('<div class="calendarNext"></div>');
		self.yearTitle = $('<span class="calendarYearTitle" title="双击标题全选">'+self.yy+'</span>');
		self.calendarYearMain.appendTo(self.container);
		self.prevBtn.appendTo(self.calendarYearMain);
		self.nextBtn.appendTo(self.calendarYearMain);
		self.yearTitle.appendTo(self.calendarYearMain);
		self.calendarGrid = $('<table class="calendarGrid" width="100%" border="0" cellspacing="0" cellpadding="0" year="" style=\'position:relative;z-index:1;\'></table>');
		self.calendarGrid.appendTo(self.container);
		self.dataBox = $('<div class="calendarGridDataBox" style="display:none;"></div>');
		self.dataBox.appendTo(self.container);
		self.container.find(".calendarGrid").bind("mousedown",function(e){e.preventDefault()});
	};
	self.createdCalendar();
	
	/*创建颜色样式*/
	self.colorStyle = function(data){
		self.styleCon = $('<style id="colorControl" type="text/css"></style>');
		var colorStyle = '';
		$.each(data,function(n){
			colorStyle += ('.colorControl_'+n+'{background:'+data[n].color+'}');
		});
		self.styleCon.html(colorStyle).appendTo($('head'));
	};
	self.colorStyle(self.colors);
	/*建立数组*/
	self.colorNumArr = function(){
		var colorArryNum = [];
		$.each(self.colors,function(n){
			colorArryNum.push(n);
		});
		return colorArryNum;
	};
	/*颜色控制器*/
	self.colorControl = function(data){
		self.control.empty();
		var colorStyle = '';
		if(data){
			$.each(data,function(n){
				var colorNum = self.getColor(data[n].code);
				var colorItem = $('<span class="colorItem colorControl_'+colorNum+'" colorCls="colorControl_'+colorNum+'" mark ="'+data[n].code+'" id="'+data[n].id+'">'+data[n].periodName+'</span>');
				colorItem.appendTo(self.control);
			});
		}
		
		self.editBtn = $('<span class="editColorBtn">编辑</span>').appendTo(self.control);
		self.editBtn.bind("click",function(){
			self.editCallback();
		});
		self.choseMsg = $('<span class="choseTxt">当前选中：<span class="markTxt"></span></span>').appendTo(self.control);
		var choseColor = '';
		self.control.find(".colorItem").bind("click",function(e){
			e.preventDefault();
			if(!$(this).hasClass("active")){
				$(this).addClass("active").siblings().removeClass("active");
				choseColor = $(this).attr('colorCls');
				self.markDate(choseColor);
				self.choseMsg.show().find(".markTxt").text($(this).text());
			}else{
				$(this).removeClass("active");
				choseColor = '';
				self.cancelMarkDate();
				self.choseMsg.hide();
			};
		});
		//alert(self.container.find(".calendarGrid").attr("year"))
		
		self.dataOption = data;
		/*移除多余标记*/
		var markArr = self.colorNumArr();
		for(j = 0;j < self.control.find(".colorItem").length;j ++){
			 for (var i = 0; i < markArr.length; i++) {
				if (markArr[i] == self.getColor(self.control.find(".colorItem").eq(j).attr("mark"))){
					markArr.splice(i, 1);
				}
			}
		};
		$.each(markArr,function(n){
			self.container.find(".calendarGrid tr").find("td.colorControl_"+markArr[n]+"").removeClass().attr("colorcls","");
		})
	};
	if(self.controlShow){
		self.colorControl(self.dataOption);
	};
	/*选择标记日期*/
	self.beforeColor = '';
	self.markDate = function(colorCls){
		self.beforeColor = colorCls;
		/*点击，拖动日期表标记选中日期*/
		self.cancelMarkDate();
		self.container.find(".calendarGrid td").bind("mousedown",function(e){
			e.preventDefault();
			if($(this).attr("hasDay") && !$(this).hasClass('noChoses')){
				if($(this).attr("colorcls") == ""){
					$(this).addClass(colorCls).addClass("active").attr("colorcls",colorCls);
				}else if($(this).attr("colorcls") != "" && $(this).attr("colorcls") != colorCls){
					$(this).removeClass($(this).attr("colorcls")).addClass(colorCls).addClass("active").attr("colorcls",colorCls);
				}else if($(this).attr("colorcls") == colorCls){
					$(this).removeClass(colorCls+" active").attr("colorcls","");
				};
			};
			self.container.find(".calendarGrid td").bind("mouseenter",function(){
				if($(this).attr("hasDay") && !$(this).hasClass('noChoses')){
					if($(this).attr("colorcls") == ""){
						$(this).addClass(colorCls).addClass("active").attr("colorcls",colorCls);
					}else if($(this).attr("colorcls") != "" && $(this).attr("colorcls") != colorCls){
						$(this).removeClass($(this).attr("colorcls")).addClass(colorCls).addClass("active").attr("colorcls",colorCls);
					}else if($(this).attr("colorcls") == colorCls){
						$(this).removeClass(colorCls+" active").attr("colorcls","");
					};
				};
			});
		});
		
		//鼠标放置在红色背景触发事件
		self.container.find(".calendarGrid td").bind("mouseover",function(e){
			var class1 = $(this).attr("class");
			var value = $(this).html();
			if(value != ""){
				if(class1 == 'colorControl_0 active'){
					$(this).attr("title","不可使用");
				}
			}
		});
		
		/*
		self.container.find(".calendarGrid tr").eq(0).bind("mouseup",function(e){
			$(this).prev().attr("<thead data-options='frozen:true'>");
			$(this).next().attr("</thead>");
		});*/
		
		/*松开鼠标按钮，注销标记动作*/
		$(document).bind("mouseup.calendarGridMouseup",function(){
			self.getYearValue();
			self.container.find(".calendarGrid td").unbind("mouseenter");
		});
		/*双击星期，标记该星期下所有日期*/
		self.container.find(".calendarGrid .dayTh").bind("dblclick",function(e){
			e.preventDefault();
			if(!$(this).hasClass("active")){
				$(this).addClass("active");
				for(var i = 0; i <= 12;i++){
					if($(this).parent().find("td").eq(i).attr("hasDay") && !$(this).parent().find("td").eq(i).hasClass('noChoses')){
						$(this).parent().find("td").eq(i).removeClass().addClass(colorCls+" active").attr("colorcls",colorCls);
					};
				};
			}else{
				$(this).removeClass("active");
				//if(){
					//$(this).parent().find("td").removeClass().attr("colorcls","");
				//}
				for(var i = 0; i <= 12;i++){
					if($(this).parent().find("td").eq(i).attr("hasDay") && !$(this).parent().find("td").eq(i).hasClass('noChoses')){
						$(this).parent().find("td").eq(i).removeClass().attr("colorcls","");
					};
				};
			};
			self.getYearValue();
		});
		/*双击月份，标记该月份下所有日期*/
		self.container.find(".calendarGrid .monthTh").bind("dblclick",function(e){
			e.preventDefault();
			var numz = $(this).index()-1;
			if(!$(this).hasClass("active")){
				$(this).addClass("active");
				for(var i = 1; i < 38; i ++){
					if(self.container.find(".calendarGrid tr").eq(i).find("td").eq(numz).attr("hasDay") && !self.container.find(".calendarGrid tr").eq(i).find("td").eq(numz).hasClass('noChoses')){
						self.container.find(".calendarGrid tr").eq(i).find("td").eq(numz).removeClass().addClass(colorCls+" active").attr("colorcls",colorCls);
					};
				};
			}else{
				$(this).removeClass("active");
				for(var i = 1; i < 38; i ++){
					if(self.container.find(".calendarGrid tr").eq(i).find("td").eq(numz).attr("hasDay") && !self.container.find(".calendarGrid tr").eq(i).find("td").eq(numz).hasClass('noChoses')){
						self.container.find(".calendarGrid tr").eq(i).find("td").eq(numz).removeClass().attr("colorcls","");
					};
				};
			};
			self.getYearValue();
		});
		self.yearTitle.bind("dblclick",function(e){
			e.preventDefault();
			if(!$(this).hasClass("active")){
				$(this).addClass("active");
				for(var i = 1; i < 38; i ++){
					for(var h = 0;h<13;h++){
						if(self.container.find(".calendarGrid tr").eq(i).find("td").eq(h).attr("hasDay") && !self.container.find(".calendarGrid tr").eq(i).find("td").eq(h).hasClass('noChoses')){
							self.container.find(".calendarGrid tr").eq(i).find("td").eq(h).removeClass().addClass(colorCls+" active").attr("colorcls",colorCls);
						};
					}
				};
			}else{
				$(this).removeClass("active");
				for(var i = 1; i < 38; i ++){
					for(var h = 0;h<13;h++){
						if(self.container.find(".calendarGrid tr").eq(i).find("td").eq(h).attr("hasDay") && !self.container.find(".calendarGrid tr").eq(i).find("td").eq(h).hasClass('noChoses')){
							self.container.find(".calendarGrid tr").eq(i).find("td").eq(h).removeClass().attr("colorcls","");
						};
					};
				};
			};
			self.getYearValue();
		});
	};
	/*注销标记日期动作*/
	self.cancelMarkDate = function(){
		self.container.find(".calendarGrid td").unbind("mousedown");
		self.container.find(".calendarGrid td").unbind("mouseenter");
		self.container.find(".calendarGrid .dayTh").unbind("dblclick");
		self.container.find(".calendarGrid .monthTh").unbind("dblclick");
		$(document).unbind("mouseup.calendarGridMouseup");
	};
	//获得一个月的第一天是星期几
	self.getFirstWeek = function(year, month) {
		var date = new Date(year, month - 1, 1);
		return date.getDay();
	};
	/*获取当前月份总天数*/
	self.getmaxDay = function(year, month) {
		var date = new Date(year, month, 0);
		return date.getDate();
	};
	
	/*初始化日期*/
	self.initDayNumbers =function(year, month) {
		var maxday = self.getmaxDay(year, month),
			startIndex = self.getFirstWeek(year, month);
		var noChosesAble = false,noChoses = '';
		var mxD = 0;
		if(year == _date.year){
			if(month < _date.month){
				noChosesAble = true;
				mxD = maxday + 1;
			};
			if(month == _date.month){
				noChosesAble = true;
				mxD = _date.day
			};
		};
		for (var s = 1; s <= maxday; s ++) {
			if(noChosesAble && mxD > s){noChoses = 'noChoses'}else{noChoses = ''};
			self.container.find(".calendarGrid tr").eq(startIndex + 1).find("td").eq(month - 1).text(s).attr('data-date', self.dayNum).attr("hasDay",true).attr('colorCls','').addClass(noChoses);
			self.dayNum ++;
			startIndex ++;
		};
    };
	/*初始化标记*/
	self.loadMark = function(data){
		//if(data){
			var nowYear = true;
			for(var n = 0; n < data.length;n++){
				self.creatDataBox(data[n].year,data[n].time).appendTo(self.dataBox);
				if(data[n].year == _date.year){
					nowYear = false;
				};
			};
			if(nowYear){
				var yVal = '';
				for(var u = 0; u < self.container.find(".calendarGrid tr td[hasDay = true]").length; u ++){
					yVal += "0";
				};
				self.creatDataBox(_date.year,yVal).appendTo(self.dataBox);			
			};
		//};	
	};
	/*加载标记*/
	self.mark = function(year){
		if(year){
			var dataz = '',n;
			n = self.yearNum(year);
			if(n != -1){
				dataz = self.dataBox.find("input").eq(n).attr("value");
				var tt = 0;
				var ds= String(dataz);
				for(var n = 0; n<=ds.length;n++){
					if(ds[n] != "0"){
						tt = n + 1;
						self.container.find(".calendarGrid tr td[data-date = "+tt+"]").addClass("colorControl_"+self.getColor(ds[n])+" active").attr("colorcls","colorControl_"+self.getColor(ds[n]));
						self.dataBox.find("input").eq(n).attr("change","false");
					};
				};
			}
		};
	};
	
	/*遍历选中年份下日期*/
	self.setCalendarYear = function(year_){
		self.calendarGrid.attr("year",year_);
		self.calendarGrid.empty();
		self.createdGrid();
		self.dayNum = 1;
		for(var m = 1; m <= 12;m ++){
			self.initDayNumbers(year_,m);	
		};
	};
	self.setCalendarYear(_date.year);
	/*年份选择器*/
	self.choseYear = function(year,step){
		var yeart = year;
		yeart += step;
		self.yearTitle.text(yeart);
		self.setCalendarYear(yeart);
		self.control.find(".colorItem").removeClass("active");
		self.mark(yeart);
		self.creatDataBox(yeart,"");
		return yeart;
	};
	/*点击选择年份*/
	self.buttonChange = function(){
		self.year_ = self.yy;
		self.container.find(".calendarPrev").bind("click",function(){
			if(self.year_ - 1 >= self.minYear){
				self.year_ = self.choseYear(self.year_,-1);
				if(self.creatDataBox(self.year_,"")){
					self.creatDataBox(self.year_,"").prependTo(self.dataBox);
				};
				self.choseMsg.hide();
			};
		});
		self.container.find(".calendarNext").bind("click",function(){
			if(self.year_ + 1 <= self.maxYear){
				self.year_ = self.choseYear(self.year_,1);
				if(self.creatDataBox(self.year_,"")){
					self.creatDataBox(self.year_,"").appendTo(self.dataBox);
				};
				self.choseMsg.hide();
			}
			
		});
	};
	self.buttonChange();
	if(self.markData.length != 0){
		self.loadMark(self.markData);
		self.mark(self.yy);
	}else{
		self.creatDataBox(self.yearTitle.text(),"").appendTo(self.dataBox);	
	};
	
	//
	
	/*获取当前年份选中值*/
	self.getYearValue = function(){
		var colorCode ="",zz = 0;
		for(var t = 0; t <= 13; t ++){
			for(var i = 0;i < 38; i++){
				if(self.container.find("tr").eq(i).find("td").eq(t).attr("data-date")){
					if(self.container.find(".calendarGrid tr").eq(i).find("td").eq(t).hasClass("active")){
						zz = self.container.find("tr").eq(i).find("td").eq(t).attr("colorcls").split("_")[1];
						colorCode += self.colors[zz].code;
					}else{
						colorCode += 0;
					};
				};
			};
		};
		self.dataBox.find("input").eq(self.yearNum(self.year_)).val(colorCode);
		if(colorCode){
			self.dataBox.find("input").eq(self.yearNum(self.year_)).attr("change","true");
		}
	};
	/*获取选中值*/
	self.getValue = function(){
		var arr = [],ii = 0;
		for(var i = 0; i < self.dataBox.find("input").length;i++){
//			if(self.dataBox.find("input").eq(i).attr("change") == "true"){
				arr[ii]= {};
				arr[ii]['year'] = self.dataBox.find("input").eq(i).attr("year");
				arr[ii]['time'] = self.dataBox.find("input").eq(i).attr("value");
				ii ++
//			};
		};
		return arr;
	};
	
	/*获取色板*/
	self.getColorValue = function(){
		return self.dataOption;
	};
}
