// JavaScript Document
function printTemplate(option){
	var self = this;
	option = $.extend({
		tableContent:null,
		container : null,
		editDetail:null,
		canvasWidth : 100,
		canvasHeight : 100,
		canvasBg:"#f6f8f7",
		data : null,
		tableClickCallback:function(type,id){},
		onmove:function(x,y){}
	},option||{});	
	
	/*初始化参数*/
	self.init = function(){
		self.proportion = 3.78;
		self.canvas = document.createElement('canvas');
		/*if($.browser.msie && $.browser.version < 9){
			self.canvas = window.G_vmlCanvasManager.initElement(self.canvas);
		};*/
		self.container = option.container,
		self.tableContent = option.tableContent,
		self.editDetail = option.editDetail,
		self.canvasWidth = option.canvasWidth*self.proportion,
		self.canvasHeight = option.canvasHeight*self.proportion,
		self.data = option.data,
		self.tableClickCallback = option.tableClickCallback;
		self.onmove = option.onmove;
		self.choseList = [];
		self.container.empty();
		self.canvasSize = '';
		self.canvasWidth >= self.canvasHeight?self.canvasSize = self.canvasWidth:self.canvasSize = self.canvasHeight;
		self.container.css({
			"width":self.canvasWidth,
			"height":self.canvasHeight,
			"overflow":"hidden",
			"position":"relative"	
			
		})
		self.canvas.width = self.canvasWidth;
		self.canvas.height = self.canvasHeight;
		self.canvas.style.position = "absolute";
		self.canvas.style.zIndex = "0";
		self.canvas.style.cursor = "pointer";
		self.container.append(self.canvas);
		$("<div class='printMain'></div>").css({
			"width":self.canvasWidth,
			"height":self.canvasHeight,
			"position":"absolute",
			"top":"0",
			"left":"0",
			"background":option.canvasBg
		}).appendTo(self.container);
	
		self.context = self.canvas.getContext("2d");
		self.choseNum = '';
		self.moveStart = false;
		self.beginX = 0,
		self.beginY = 0;
		self.maxLeft = 0;
		self.maxTop = 0;
		
	};
	self.init();
	/*计算是否为新增内容*/
	self.getNewIdNum = function(id){
		var t = -1;
		for(var g = 0; g < self.choseList.length;g++){
			if(self.choseList[g].id == id){
				t = g;
				break;
			};
		};
		return t;
	};
	/*画图*/
	self.printTicket = function(data){
		
		self.container.find(".printMain").empty();
		$.each(data,function(n){
			var textDiv = $('<div class="moveText"></div>').css({
					"transform":"rotate("+data[n].orientation+"deg)",
					"-ms-transform":"rotate("+data[n].orientation+"deg)",	
					"-moz-transform":"rotate("+data[n].orientation+"deg)", 	
					"-webkit-transform":"rotate("+data[n].orientation+"deg)", 
					"-o-transform":"rotate("+data[n].orientation+"deg)",
					"left":data[n].seatX*self.proportion+"px",
					"top": data[n].seatY*self.proportion+"px",
					"position":"absolute",
					"white-space": "nowrap",
					"display":"inline-block"
				});
			if(data[n].font == "QR"){
				textDiv.css({
					"width":data[n].fontSizeRate*data[n].width*self.proportion,
					"height":data[n].fontSizeRate*data[n].height*self.proportion
				}).html("<img src='/systemUI/3.0/images/codeImg/qrcode.png' width='"+data[n].fontSizeRate*data[n].width*self.proportion+"' height='"+data[n].fontSizeRate*data[n].height*self.proportion+"'/>");
			}else if(data[n].font == "Code128"){
				textDiv.css({
					"width":data[n].fontSizeRate*data[n].content.length*self.proportion,
					"height":data[n].barcodeHeight*self.proportion
				}).html("<img src='/systemUI/3.0/images/codeImg/barcode.png' width='"+data[n].fontSizeRate*data[n].content.length*self.proportion+"' height='"+data[n].barcodeHeight*self.proportion+"'/>");
			}else{
				textDiv.css({
					"font-size":data[n].fontSizeRate+"pt"
				}).html(data[n].content);
			}
			
			textDiv.appendTo(self.container.find(".printMain"));
			if(self.choseNum.length != 0 && n == self.choseNum){textDiv.addClass("focus")};
			var textWidth = textDiv.width(),textHeight = textDiv.height();
			var origin = "50% 50%";
				self.marginT = 0,self.marginL = 0;
			switch(data[n].orientation){
				case "90":
					//origin = "0% 100%";
					self.marginL = (textWidth-textHeight)/2;
					self.marginT = (textHeight-textWidth)/2;
					break;
				/*case "180":
					origin = "100% 100%";
					self.marginL = -textWidth;
					self.marginT = -textHeight;
					break;*/
				case "270":
					//origin = "100% 0%";
					self.marginL = (textWidth-textHeight)/2;
					self.marginT = (textHeight-textWidth)/2;
					break;
			};
			textDiv.css({
				"width":textWidth,
				"height":textHeight,
				"transform-origin":origin,
				"-ms-transform-origin":origin,
				"-moz-transform-origin":origin,
				"-webkit-transform-origin":origin,
				"-o-transform-origin":origin,
				"margin-top":-self.marginT +"px",
				"margin-left":-self.marginL +"px"
			});
		});
		self.mouseListener();
		//self.context.closePath(); 
	};
	
	/*下拉列表选中项*/
	self.selectedFunc = function(obj,selectOn){
		for(var j = 0; j < obj.length;j ++){
			if(obj.eq(j).val() == selectOn){
				obj.eq(j).attr("selected","selected");
				break;	
			};
		};	
	};
	/*加载选中样式数据列*/
	self.choseTr = function(n){
		var dataMenu = self.choseList[n];
		if(dataMenu){
			self.editDetail.find("select option").removeAttr("selected");
			self.editDetail.find("input[name = seatX]").val(dataMenu.seatX);
			self.editDetail.find("input[name = seatY]").val(dataMenu.seatY);
			/*字体选项*/
			self.selectedFunc(self.editDetail.find("select[name = font] option"),dataMenu.font);
			/*字号选项*/
			self.selectedFunc(self.editDetail.find("select[name = fontSize] option"),dataMenu.fontSize);
			/*方向选项*/
			self.selectedFunc(self.editDetail.find("select[name = orientation] option"),dataMenu.orientation);
			/***/
			self.editPanel($("#elementTab tr").eq(n+1).addClass("focus"));
			
		}else{
			self.uneditPanel();	
		};
	};
	//seatX:"12",seatY:"12",font:"Arial",fontSize:"12",orientation:"0"
	/*修改表格数据*/
	self.tableData = function(obj,field,val){
		$(obj).find("td[field="+field+"]").text(val);
	};
	
	/*编辑样式*/
	self.editPanel = function(){
		var editObj = null;
		
		self.editDetail.find("input[type = text]").bind("keyup",function(){
			switch($(this).attr("name")){
				case "seatX":
					self.renewPositionX(self.choseNum,$(this).val());
					break;
				case "seatY"	:
					self.renewPositionY(self.choseNum,$(this).val());
					break;
			}
		});
		self.editDetail.find("select").bind("change",function(){
			
			switch($(this).attr("name")){
				case "font":
					self.renewFontFamily(self.choseNum,$(this).val());
					break;
				case "fontSize"	:
					//self.renewFontSize(self.choseNum,$(this).val());
					break;
				case "orientation" :
					self.renewFontDirection(self.choseNum,$(this).val());
					break;	
			}
		});
	};
	//self.editPanel();
	/*注销编辑功能*/
	self.uneditPanel = function(){
		self.editDetail.find("input[type = text]").unbind("keyup");
		self.editDetail.find("select").unbind("change");
	};
	/*更新位置X*/
	self.renewPositionX = function(id,x){
		var index = self.getNewIdNum(id);
		x = x*self.proportion;
		setMaxActiveArea(index);
		if(self.maxLeft >= x){
			self.choseList[index].seatX = x/self.proportion;						
			$("input[name = seatX]").val(roundNum(x/self.proportion));
			$("#elementTab tr[id='"+id+"']").find("td").eq(2).text(roundNum(x/self.proportion));
		}else{
			$("input[name = seatX]").val(roundNum(self.maxLeft/self.proportion));
			$("#elementTab tr[id='"+id+"']").find("td").eq(2).text(roundNum(self.maxLeft/self.proportion));
			self.choseList[index].seatX = self.maxLeft/self.proportion;						
		};
		
		self.printTicket(self.choseList);
	};
	/*更新位置Y*/
	self.renewPositionY = function(id,y){
		var index = self.getNewIdNum(id);
		setMaxActiveArea(index);
		y = y*self.proportion;
		if(self.maxTop >= y){
			self.choseList[index].seatY = y/self.proportion;
			self.printTicket(self.choseList);
			$("input[name = seatY]").val(roundNum(y/self.proportion));
			$("#elementTab tr[id='"+id+"']").find("td").eq(3).text(roundNum(y/self.proportion));
		}else{
			$("input[name = seatY]").val(roundNum(self.maxTop/self.proportion));
			$("#elementTab tr[id='"+id+"']").find("td").eq(3).text(roundNum(self.maxTop/self.proportion));
			self.choseList[index].seatY = self.maxTop/self.proportion;
			self.printTicket(self.choseList);
		}
	};
	/*更新字体大小*/
	self.renewFontSize = function(id,size){
		var index = self.getNewIdNum(id);
		self.choseList[index].fontSizeRate = size;
		self.printTicket(self.choseList);
		setMaxActiveArea(index);
		self.renewPositionX(id,$("input[name = seatX]").val());
		self.renewPositionY(id,$("input[name = seatY]").val());
	};
	/*更新字体*/
	self.renewFontFamily = function(id,fontFamily,rate,w,h,barh){
		var index = self.getNewIdNum(id);
		self.choseList[index].font = fontFamily;
		self.choseList[index].fontSizeRate = rate;
		self.choseList[index].width = w;
		self.choseList[index].height = h;
		self.choseList[index].barcodeHeight = barh;
		self.printTicket(self.choseList);
		setMaxActiveArea(index);
		self.renewPositionX(id,$("input[name = seatX]").val());
		self.renewPositionY(id,$("input[name = seatY]").val());
	};
	/*更新字体角度*/
	self.renewFontDirection = function(id,direction){
		var index = self.getNewIdNum(id);
		self.choseList[index].orientation = direction;
		self.printTicket(self.choseList);
		setMaxActiveArea(index);
		self.renewPositionX(id,$("input[name = seatX]").val());
		self.renewPositionY(id,$("input[name = seatY]").val());
	};
	/*写入内容*/
	self.newContent = function(id,content){		
		var index = self.getNewIdNum(id);		
		self.choseList[index].content = content;
		self.printTicket(self.choseList);
	};
	/*改变code大小*/
	self.codeSize = function(id,size){
		var index = self.getNewIdNum(id);
		self.choseList[index].barcodeHeight = size;
		
		self.printTicket(self.choseList);
	}
	/*改变画布尺寸*/
	/*宽度*/
	self.resizeCanvasW = function(w){
		w*self.proportion >= self.canvasHeight?self.canvasSize = w*self.proportion:self.canvasSize = self.canvasHeight;
		self.container.css({
			"width":w*self.proportion
		});
		self.container.find(".printMain").width(w*self.proportion);
		self.canvas.width = w*self.proportion;
		self.canvasWidth = w*self.proportion;
		$.each(self.choseList,function(n){
			self.renewPositionX(self.choseList[n].id,self.choseList[n].seatX);		
			self.onmove(roundNum(self.choseList[n].seatX),roundNum(self.choseList[n].seatY),$("#elementTab tr").eq(n+1).attr("id"));
			$("input[name = seatX]").val($("#elementTab tr.focus td[field=seatX]").text());
		});
	};
	/*高度*/
	self.resizeCanvasH = function(h){
		self.canvasWidth >= h*self.proportion?self.canvasSize = self.canvasWidth:self.canvasSize = h*self.proportion;
		self.container.css({
			"height":h*self.proportion
		});
		self.container.find(".printMain").height(h*self.proportion);
		self.canvas.height = h*self.proportion;
		self.canvasHeight = h*self.proportion;
		$.each(self.choseList,function(n){
			self.renewPositionY(self.choseList[n].id,self.choseList[n].seatY);
			self.onmove(roundNum(self.choseList[n].seatX),roundNum(self.choseList[n].seatY),$("#elementTab tr").eq(n+1).attr("id"));
			$("input[name = seatY]").val($("#elementTab tr.focus td[field=seatY]").text());
		});
	};
	/*鼠标事件监听*/
	self.mouseListener = function(){
		var moveAble = false;
		
		self.container.find(".printMain .moveText").bind("mousedown",function(event){
			event.preventDefault();
			 moveAble = true;
			 self.moveStart = true;
			 $(this).addClass("focus").siblings().removeClass("focus");
			 self.canvas.style.zIndex = "999";
			 self.choseNum = $(this).index();
			 $("#elementTab tr").eq(self.choseNum+1).addClass("focus").siblings().removeClass("focus");
			 $("#elementTab tr").eq(self.choseNum+1).click();
			 setMaxActiveArea(self.choseNum);
			 moveMaxArea(self.choseNum)
		}); 
		self.canvas.addEventListener("mousemove", function (evt) {
			if(moveAble && self.moveStart){
			  var mousePos = getMousePos(self.canvas, evt);
			 
			  /*xu*/
				if(self.moveMinLeft<= mousePos.x && mousePos.x <= self.moveMaxLeft && self.moveMinTop<= mousePos.y && mousePos.y <=self.moveMaxTop){
					 $("#msg").val( roundNum((mousePos.x- self.offsetLeft)/self.proportion)+","+roundNum((mousePos.y- self.offsetTop)/self.proportion));
					 self.choseList[self.choseNum].seatX = (mousePos.x- self.offsetLeft)/self.proportion;
					 self.choseList[self.choseNum].seatY = (mousePos.y- self.offsetTop)/self.proportion;
					 self.container.find(".printMain .moveText.focus").css({
						 "left":mousePos.x- self.offsetLeft,
						 "top":mousePos.y- self.offsetTop
					})
				  // self.renewPositionX(self.choseNum,mousePos.x);
				  // self.renewPositionY(self.choseNum,mousePos.y);
					$("input[name = seatX]").val(roundNum((mousePos.x- self.offsetLeft)/self.proportion));
					$("input[name = seatY]").val(roundNum((mousePos.y- self.offsetTop)/self.proportion));
				   self.onmove(roundNum((mousePos.x- self.offsetLeft)/self.proportion),roundNum((mousePos.y- self.offsetTop)/self.proportion),$("#elementTab tr").eq(self.choseNum+1).attr("id"));
			
				}
			 
			}
		}, false); 
		$(document).bind('mouseup', function (event) { 
				moveAble = false;
				self.moveStart = false;
				self.canvas.style.zIndex = "0";
		}); 
		self.canvas.addEventListener('mouseup', function (event) { 
			  moveAble = false;
			  self.moveStart = false;
			  self.canvas.style.zIndex = "0";
		}, false); 
			
		};
/*文字最大活动区域计算*/
	function setMaxActiveArea(index){
		
		var textW = self.container.find(".printMain .moveText").eq(index).width(),
	  	  textH = self.container.find(".printMain .moveText").eq(index).height();
		 
		  switch(self.choseList[index].orientation){
			case "0":
				self.maxLeft = self.canvasWidth - textW;
				self.maxTop = self.canvasHeight - textH;
				break;
			case "90":
				self.maxLeft = self.canvasWidth - textH;
				self.maxTop = self.canvasHeight - textW;
				break;
			case "180":
				self.maxLeft = self.canvasWidth - textW;
				self.maxTop = self.canvasHeight - textH;
				break;
			case "270":
				self.maxLeft = self.canvasWidth - textH;
				self.maxTop = self.canvasHeight - textW;
				break;
			};
			self.maxLeft<0?self.maxLeft = 0:self.maxLeft = self.maxLeft;
			self.maxTop<0?self.maxTop = 0:self.maxTop = self.maxTop;
	}
	/*鼠标拖动最大范围*/
	function moveMaxArea(index){
		var textW = self.container.find(".printMain .moveText").eq(index).width(),
	  	  textH = self.container.find(".printMain .moveText").eq(index).height();
			 self.offsetLeft = textW/2;
			 self.offsetTop = textH/2;
		  switch(self.choseList[index].orientation){
			case "0":
				self.moveMinLeft = textW/2;
				self.moveMinTop = textH/2;
				self.moveMaxLeft = self.canvasWidth - textW/2;
				self.moveMaxTop = self.canvasHeight - textH/2;				
				break;
			case "90":
				self.moveMinLeft = textH/2;
				self.moveMinTop = textW/2;
				self.moveMaxLeft = self.canvasWidth - textH/2;
				self.moveMaxTop = self.canvasHeight - textW/2;
				self.offsetLeft = textH/2;
				self.offsetTop = textW/2;
				break;
			case "180":
				self.moveMinLeft = textW/2;
				self.moveMinTop = textH/2;
				self.moveMaxLeft = self.canvasWidth - textW/2;
				self.moveMaxTop = self.canvasHeight - textH/2;
				break;
			case "270":
				self.moveMinLeft = textH/2;
				self.moveMinTop = textW/2;
				self.moveMaxLeft = self.canvasWidth - textH/2;
				self.moveMaxTop = self.canvasHeight - textW/2;
				self.offsetLeft = textH/2;
				self.offsetTop = textW/2;
				break;
			};
	
	};
	/*获取鼠标坐标*/
	function getMousePos(canvas, evt) {
	   var rect = canvas.getBoundingClientRect();
	   return {
		 x: evt.clientX - rect.left * (canvas.width / rect.width),
		 y: evt.clientY - rect.top * (canvas.height / rect.height)
	   }
	 }
	self.mouseListener();
	function roundNum(num){
		return Math.round(num);		
	};
	/*绑定表格点击事件*/
	self.tableChose = function(){
		self.tableContent.find("tr").unbind("mouseup.tableClick");
		self.tableContent.find("tr").bind("mouseup.tableClick",function(){
			$(this).addClass("focus").siblings().removeClass("focus");
			self.choseNum = $(this).index();
			self.moveStart = true;
			self.container.find(".printMain .moveText").eq(self.choseNum).addClass("focus").siblings().removeClass("focus");
		});
	};
	/*设置选中项*/
	self.focusIndex = function(id){
		var index = self.getNewIdNum(id);
		self.container.find(".printMain .moveText").eq(index).addClass("focus").siblings().removeClass("focus");
		self.choseNum = index;
		//self.setMaxActiveArea(self.choseNum);
	};
	/*设置index*/
	self.setIndex = function(id){
		var index = self.getNewIdNum(id);
		self.choseNum = index;
		self.moveStart = true;
	};
	/*载入数据*/
	self.loadData = function(data){
		if(self.getNewIdNum(data[0].id) == -1){
			self.choseList[self.choseList.length] = data[0];
		};
		self.printTicket(self.choseList);
		//self.tableChose();
		self.mouseListener();
	};
	/*删除行*/
	self.deleteRow = function(id){
		var index = self.getNewIdNum(id);
		if(index != -1){
			self.choseList.splice(index,1)
			self.printTicket(self.choseList);
		}
		
	}
}