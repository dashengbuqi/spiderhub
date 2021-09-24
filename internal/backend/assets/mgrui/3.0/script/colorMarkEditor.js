// JavaScript Document

function colorMarkEditor(option){
	var self = this;
	option = $.extend({
		containerz:null,
		dataColorz:null,
		controlBoxColor:null
	},option||{});
	/**/
	self.initz = function(){
		self.containerz = option.containerz,
		self.dataColorz = option.dataColorz,
		self.controlBoxColor = option.controlBoxColor;
		self.arrColor = self.controlBoxColor; 
		self.datalength = self.dataColorz.length;
	};
	self.initz();
	self.newData = [],self.deleteData = [],self.newDataNum = 0,self.deleteDataNum = 0;
	/*更新颜色列表*/
	self.renewColorList = function(datac){
		//debugger;
		self.colorMenu.find("ul").empty();
		$.each(self.controlBoxColor,function(n){
			var lis = $('<li style="background:'+self.controlBoxColor[n].color+'" mark="'+self.controlBoxColor[n].code+'"></li>');
			lis.appendTo(self.colorMenu.find("ul"));
		});
		$.each(datac,function(n){
			for(var y = 0; y < self.controlBoxColor.length;y++){
				if(self.controlBoxColor[y].color == self.getColorz(datac[n].code)){
					self.colorMenu.find("ul li[mark = "+self.controlBoxColor[y].code+"]").remove();
					break;
				};
			};
		});
		self.colorEditor.find(".colorMark").css("background",self.getColorz(self.colorMenu.find("ul li:first-child").attr("mark"))).attr("markc",self.colorMenu.find("ul li:first-child").attr("mark"))
	};
	
	/**/
	self.getColorNumz = function(code){
		for(var g = 0; g < self.controlBoxColor.length;g++){
			if(self.controlBoxColor[g].code == code){
				return g;
				break;
			};
		};
	};
	self.getColorz = function(code){
		for(var g = 0; g < self.controlBoxColor.length;g++){
			if(self.controlBoxColor[g].code == code){
				return self.controlBoxColor[g].color;
				break;
			};
		};
	};
	/*计算新增颜色位置*/
	self.getNewColorNum = function(code){
		//alert(self.newData.length)
		for(var g = 0; g < self.newData.length;g++){
			if(self.newData[g].code == code){
				return g;
				break;
			};
		};
	};
	self.menuListClick=function(){
		self.colorMenu.find("li").bind("click",function(){
			self.colorEditor.find(".colorMark").css("background",self.getColorz($(this).attr("mark")));
			self.colorEditor.find(".colorMark").attr("markc",$(this).attr("mark"));
			self.colorMenu.hide();
		});
	};
	/**/
	self.deleteItem = function(){
		self.colorList.find("table .deleteColor").unbind("click");
		self.colorList.find("table .deleteColor").bind("click",function(){
			$(this).parent().parent().remove();
			self.renewColorList(self.getValuez());
			self.menuListClick();
			//console.log(self.dataColorz[self.getColorNumz($(this).parent().parent().attr("mark"))])
			if(self.dataColorz[self.getColorNumz($(this).parent().parent().attr("mark"))]){
				self.deleteData[self.deleteDataNum]= {};
				self.deleteData[self.deleteDataNum]['id'] = self.dataColorz[self.getColorNumz($(this).parent().parent().attr("mark"))].id;
				self.deleteDataNum ++;
			};
			//alert($(this).parent().parent().attr("mark"))
			
			//alert(self.getNewColorNum($(this).parent().parent().attr("mark")))
			//console.log(self.getNewColorNum($(this).parent().parent().attr("mark")))
			
			var st = self.getNewColorNum($(this).parent().parent().attr("mark"));
			if(st){
				self.newData.splice(st,1);
				self.newDataNum --;
			};
			
			//
		});
	};
	/*加载表格数据*/
	self.createdGrid = function(datas){
		self.colorList.find("table").empty();
		self.tableHeader = $('<tr><th scope="col">标记颜色</th><th scope="col">标记名称</th><th scope="col">操作</th></tr>');
		self.tableHeader.appendTo(self.colorList.find("table"));
		$.each(datas,function(n){
			var colorItem = $('<tr mark="'+datas[n].code+'"><td><span class="colorItemIcon" style="background:'+self.getColorz(datas[n].code)+'"></span></td><td><span class="colorItemName"><input type="text" value="'+datas[n].periodName+'" class="dataGName"/></span></td><td><span class="deleteColor">删除</span></td></tr>');
			colorItem.appendTo(self.colorList.find("table"));
		});
		self.renewColorList(datas);
		self.deleteItem();
	};
	/**/
	self.created = function(){
		self.containerz.empty();
		self.colorEditor = $('<div class="colorEditorMain"><div class="colorEditorLabel">新增时段：</div><div class="colorEditor"><span class="colorMark" markc="'+self.controlBoxColor[2].code+'"></span><input type="text" id="textValue"></div><a href="#" class="borderBlueBtn" onclick="">增加</a></div>');
		self.colorEditor.appendTo(self.containerz);
		self.colorMenu = $('<div class="colorMenus"><ul></ul></div>');
		self.colorMenu.appendTo(self.colorEditor);
		
		self.colorList = $('<div class="colorListMain"><table class="colorGrid" width="100%" border="0" cellspacing="0" cellpadding="0"></table></div>');
		self.colorList.appendTo(self.containerz);
		self.createdGrid(self.dataColorz);
		self.renewColorList(self.dataColorz);
	};
	self.created();
	
	/**/
	self.clickColorMark = function(){
		self.colorEditor.find(".colorMark").bind("click",function(){
			self.colorMenu.toggle();
			$(document).bind("click.clickColorMark",function(e){
				if($(e.target).closest(".colorMark").length == 0 && $(e.target).closest(".colorMenus li").length == 0 && $(e.target).closest(".colorMenus li").length == 0){
					self.colorMenu.hide();
				}
			});
		});
		
	};
	
	self.menuListClick()
	self.clickColorMark();
	self.submitBtn = function(){
		self.colorEditor.find(".borderBlueBtn").bind("click",function(){
			if($("#textValue").val().length>10){
				showMsg("字段名称不能超过20位", false, false, "warning");
				$("#textValue").val("");
				return false;
			};
			if(!self.colorEditor.find(".colorEditor input[type=text]").val()){
				showMsg("请输入标记名称", false, false, "warning");
				return false;
			};
			if(self.colorMenu.find("ul li").length == 0){
				showMsg("最多能添加10种标记类型",false,false,'question');
				return false;			
			};
			var colorItem = $('<tr mark="'+self.colorEditor.find(".colorMark").attr("markc")+'"><td><span class="colorItemIcon" style="background:'+self.getColorz(self.colorEditor.find(".colorMark").attr("markc"))+'"></span></td><td><span class="colorItemName"><input type="text" value="'+self.colorEditor.find(".colorEditor input[type=text]").val()+'" class="dataGName"/></span></td><td><span class="deleteColor">删除</span></td></tr>');
			colorItem.appendTo(self.colorList.find("table"));
			self.deleteItem();
			self.newData[self.newDataNum]= {};
			self.newData[self.newDataNum]['id'] = "";
			self.newData[self.newDataNum]['periodName'] = self.colorEditor.find(".colorEditor input[type=text]").val();
			self.newData[self.newDataNum]['code'] = self.colorEditor.find(".colorMark").attr("markc");
			self.newData[self.newDataNum]['color'] = self.getColorz(self.colorEditor.find(".colorMark").attr("markc"));
			self.newDataNum ++;

				self.colorEditor.find(".colorEditor input[type=text]").val("");
				self.renewColorList(self.getValuez());
				self.menuListClick();
			
		});
	};
	self.submitBtn();
	
	self.deleteItem();
	/*点击编辑*/
	self.gridTdTxt = function(){
		self.colorList.find("table tr td .dataGName").bind("click",function(){
			//alert(0)
			self.colorList.find("table tr td .dataGName").removeClass("focus");
			var obj = $(this);
			obj.addClass("focus");
			$(document).bind("click.inputFocus",function(e){
				if($(e.target).closest(".dataGName").length == 0){
					obj.removeClass("focus");
					$(document).unbind("click.inputFocus");
				}
			})
		});
	};
	self.gridTdTxt();
	/*获取值*/
	self.getValuez = function(){
		var arr = [];
		for(var i = 1; i < self.colorList.find("table tr").length;i++){
			arr[i-1]= {};
			arr[i-1]['periodName'] = self.colorList.find("table tr").eq(i).find(".colorItemName input").val();
			arr[i-1]['code'] = self.colorList.find("table tr").eq(i).attr("mark");
			arr[i-1]['color'] = self.getColorz(self.colorList.find("table tr").eq(i).attr("mark"));
		};
		return arr;
	};
	/*获取新增标记*/
	self.getNewColor = function(){
		return self.newData;
	};
	
	self.getDeleteId = function(){
		return self.deleteData;
	};
}


