/**
 * 返回小数值
 * @param floatvar 值
 * @param decimalPlace 小数位长度
 * @returns
 */
function returnDecimal (floatvar, decimalPlace) {
	if (decimalPlace<=0) { return floatvar};
	var f_x = parseFloat(floatvar);
	if ( isNaN(f_x) ) {return floatvar;}
	var num = "1";
	for (var i=0;i<decimalPlace;i++) {num += "0";}
	var f_x = Math.round(floatvar*num)/num;
	var s_x = f_x.toString();
	var pos_decimal = s_x.indexOf('.');
	if (pos_decimal < 0) {pos_decimal = s_x.length;s_x += '.';}
	while (s_x.length <= pos_decimal + decimalPlace) {s_x += '0';}
	return s_x;
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
			imgMaxSize:102400,
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
					showMsg("选择文件错误,图片类型必须是" + opts.ImgType.join("，") + "中的一种",false, false, "warning");
					this.value = "";
					return false
				}
				_self.image = new Image();
     			_self.image.src = this.files[0];
     			if(_self.image.filesize > opts.imgMaxSize){
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
var userAgent = navigator.userAgent.toLowerCase();
//Figure out what browser is being used
jQuery.browser = {
version: (userAgent.match( /.+(?:rv|it|ra|ie)[\/: ]([\d.]+)/ ) || [])[1],
safari: /webkit/.test( userAgent ),
opera: /opera/.test( userAgent ),
msie: /msie/.test( userAgent ) && !/opera/.test( userAgent ),
mozilla: /mozilla/.test( userAgent ) && !/(compatible|webkit)/.test( userAgent )
}; 
/*获取图片地址*/
function getObjectURL(fileVal,imgID) {
	var url = null;
	if (window.createObjectURL != undefined) {
		url = window.createObjectURL(fileVal)
	} else if (window.URL != undefined) {
		url = window.URL.createObjectURL(fileVal)
	} else if (window.webkitURL != undefined) {
		url = window.webkitURL.createObjectURL(fileVal)
	}
	$("#" + imgID).attr('src', url)
	//return url
};
/*检查文件格式*/
function checkFile(file,imgID){
	var _self = this;
	var ImgType =  ["gif", "jpeg", "jpg", "bmp", "png"],imgMaxSize = '102400';
	
	if (file.value) {
		if (!RegExp("\.(" + ImgType.join("|") + ")$", "i").test(file.value.toLowerCase())) {
			showMsg("选择文件错误,图片类型必须是" + opts.ImgType.join("，") + "中的一种",false, false, "warning");
			file.value = "";
			return false
		}
		_self.image = new Image();
			_self.image.src = file.files[0];
		if(_self.image.filesize > imgMaxSize){
			showMsg("选择文件错误,图片大小必须小于"+ opts.imgMaxSize/1000+"k", false, false, "warning");
			return false
		};
		if ($.browser.msie) {
			try {
				getObjectURL(file.value,imgID);
				//$("#" + imgID).attr('src', _self.getObjectURL(this.files[0]))
			} catch (e) {
			var src = "";
			var obj = $("#" + imgID);
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
				'width': '100px',
				'height': '100px'
			});
			div.filters.item("DXImageTransform.Microsoft.AlphaImageLoader").src = src
		}
		} else {
			getObjectURL(file.value,imgID);
		}
	}

}
/**
 * UBB代码转换成HTML
 * @param ubbcode
 */
function ubb2html(ubbcode){

}