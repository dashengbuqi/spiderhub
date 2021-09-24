$.extend($.fn.validatebox.defaults.rules, {
    /* 实例 */
    minLength: {
        validator: function(value, param){
            return value.length >= param[0];
        },
        message: '至少{0}字符.'
    },
    
    /* 实例 */
    maxLength: {
        validator: function(value, param){
            return value.length <= param[0];
        },
        message: '不能超过{0}字符.'
    },
    
    /* 数字校验：number[min,max] */
    number: {
        validator: function(value, param){
            if(checkNUM(value) == 0) {
            	return false;
            }
            if(value < param[0] || value > param[1]) {
            	return false;
            }
            return true;
        },
        message: '请输入{0}~{1}的整数'
    },
    
    	/*是否包含特殊字符*/
    	nameValid: { 
        	validator: function(value) { 
        		var ch =/^([\u4e00-\u9fa5]+|[a-zA-Z0-9]+)$/;
        		return ch.test(value);
        		}, 
        	message:"请勿输入特殊字符"
        },
        
    /* 确认密码校验 */
    equalTo: {
		validator: function (value, param){
			return $(param[0]).val() == value; 
		},
		message: '两次密码输入不一致' 
	},
	
	/*下拉列表必填*/
	    comboVry: {
	        validator: function (value) {
	        	if(value == "请选择" || value == "" || value == null || value == "--请选择--" ){
	        		return false;
	        	}else{
	        		 return true;
	        	}
	        },
	        message: '请选择内容'
	    },
	  
	/* 密码格式校验*/
	pwdFormat:{
		validator:function(value){
			var ls = 0;
			if(value.match(/([a-z])+/)){
				ls++;
			}
			if(value.match(/([0-9])+/)){
				ls++; 
			}
			if(value.match(/[^a-zA-Z0-9]+/)){
				ls++;
			}
			if(value.match(/([A-Z])+/)){
				ls++;
			} 
			if(ls<3){
				return false;
			}
			return true;
		},
		message: '密码必须是大写字母，小写字母，数字，特殊字符中任意三个组合'
	},
	
    /* 数字校验2：isNumber */
    isNumber: {
        validator: function(value){
            if(checkNUM(value) == 0) {
            	return false;
            }
            return true;
        },
        message: '请输入数字'
    },
    
    /* 金额验证2：isMoney */
    isMoney: {
        validator: function(value){
            if(checkNUM(value) == 0) {
            	return false;
            }
            return true;
        },
        message: '请输入正确的金额'
    },

    /*校验邮箱*/
    isEmail: {
        validator: function(value){
            if(!Email(value)) {
            	return false;
            }
            return true;
        },
        message: '请输入正确的邮箱'
    },
    
    /* 金额校验：amount[min,max] */
    amount: {
    	validator: function(value, param) {
        	 if (isAmount(value) == 0) {  
        	 		return false;
        	 }
            if(value < param[0] || value > param[1]) {
            	return false;
            }
            return true;
    	},
    	message: '请输入{0}~{1}的金额，不超过两位小数'
    },
    
    /*验证手机号码组*/
    mobileGroup:{
    	validator:function(value){
    		/*var phoneRegex = /^1[3|4|5|7|8]\d{9}$/;
    		if(!isMobile(value)){
    			if(!/^\d{11}(,\d{11})*$/.test(value)){
    				return false;
    			}
    			return false;
    		}*/
    		if(!/^1[3|4|5|7|8]\d{9}(,1[3|4|5|7|8]\d{9})*$/.test(value)){
    			return false;
    		}
    		return true;
    	},
        message: '手机号码组格式不正确'
    },
    
    /* 金额校验：isAmount */
    isAmount: {
        validator: function(value){      
        	 if (isAmount(value) == 0) {  
        	 		return false;
        	 }
        	 return true;
        },
        message: '请输入金额，不超过两位小数'
    },
    
    /* 字符校验：string[min,max]*/
    string: {
        validator: function(value, param){
			var str = ATrim(value);
            if(str.length < param[0] || str.length > param[1]) {
            	return false;
            }
            return true;
        },
        message: '请输入{0}~{1}个字符'
    },
    /* 字符校验：string[min,max]*/
    stringmobile: {
        validator: function(value, param){
			var str = ATrim(value);
            if(str.length < param[0] || str.length > param[1]) {
            	return false;
            }
            return true;
        },
        message: '超出规定字符位,应改为区号加座机号'
    },
    
    /* 手机校验：mobile */
    mobile: {
        validator: function(value){
            return isMobile(value);
        },
        message: '手机格式不正确'
    },
    
    /* 金额校验：number[min,max] */
    money: {
        validator: function(value){
        	 var patrn = /^-?\d+\.{0,}\d{0,}$/;       
        	 if (!patrn.exec(value)) {  
        	 		return false;
        	 }
            return true;
        },
        message: '请正确输入金额'
    },
    /* 包含两位小数点的小数校验：floatNumber */
    floatNumber: {
        validator: function(value){
			// var regu = /(^([1-9]\d*|[0])\.\d{1,2}$|^[1-9]\d*$)/;
			var regu=/^\d+(\.\d{1,2})?$/ ;
			var re = new RegExp(regu);
			return re.test(value);
        },
        message: '请输入至多包含两位小数的数字'
    },
    /* 费率校验：floatNumber */
    floatRate: {
        validator: function(value){
			return !isRateErr(value);
        },
        message: '请输入正确的费率'
    },
    /* 手机校验：mobile */

    idCard: {
        validator: function(value){
            if(!Isidcard(value)) {
            	return false;
            }
            return true;
        },
        message: '身份证号码格式不正确'
    },
    /* 编号验证：numOrLetter  判断输入项是否由数字或字母或数字和字母组成*/
    numOrLetter:{
    	validator: function(value,param){
    		var str = ATrim(value);
    		if((str.length < param[0] || str.length > param[1]) || !isNumOrLetter(str)) {
            	return false;
            }
            /*if() {
            	return false;
            }*/
            return true;
        },
        message: '必须输入{0}~{1}位的字母或数字！'
    },
    
    /* 编号验证：letterNumValid  判断输入项是否由数字或字母或数字和字母组成*/
    letterNumValid:{
    	validator: function(value,param){
    		var str = ATrim(value);
    		if((str.length != param[0]) || !isNumOrLetter(str)) {
            	return false;
            }
            return true;
        },
        message: '请输入{0}位的字母或数字！'
    },
    
    chineseEG : {// 验证中文和字母 
        validator : function(value) { 
        	var ch = /^[\u4e00-\u9fa5-Za-z]{0,30}$/;
            return ch.test(value); 
        }, 
        message : '请输入中文或者字母' 
    },
    
    chinese : {// 验证中文
    	validator : function(value,param) { 
        	var ch = /^[\u4e00-\u9fa5]{0,30}$/;
        		if(!ch.exec(value)){
        			return false;
        		}
        		var str = ATrim(value);
                if(str.length > param[0]) {
                	return false;
                }
                return true;
        }, 
        message : '请输入中文,长度不能超过{0}个字符' 
    },
    
    goodsVerify: {
        validator: function (value, param) {
        	var patrn=/^[\u0391-\uFFE5\w]+$/;
        	if(!patrn.exec(value)){
        		return false;
        	}
    		var str = ATrim(value);
            if(str.length > param[0]) {
            	return false;
            }
            return true;
        },
        message: '只允许汉字、英文字母、数字及下划线,长度不能超过{0}个字符'
    },
    
    warnSpan: {
        validator: function (value) {
            return /^[\u0391-\uFFE5\w]+$/.test(value);
        },
        message: '只允许汉字、英文字母、数字及下划线。'
    },
    
    /* 金额校验：number[min,max] */
    limitedMoney: {
        validator: function(value, param){
        	 var patrn = /^-?\d+\.{0,}\d{0,}$/;       
        	 if (!patrn.exec(value) || ATrim(value).length < param[0] || ATrim(value).length > param[1]) {  
        	 		return false;
        	 }
            return true;
        },
        message: '请正确输入{0}~{1}位金额'
    },
    /*验证开始时间小于结束时间 */
    md: {
    	validator: function(value, param){
    		startTime2 = '';
		if($(param[0]).length != 0){
			startTime2 = $(param[0]).val();
		};
    	if(!startTime2){
    		var today = new Date();
    		var y = today.getFullYear();  
            var m = today.getMonth()+1;  
            var d = today.getDate();
            var da = (y+"-"+m+"-"+d);
    		startTime2 = da;	
    	}
    	var d1 = $.fn.datebox.defaults.parser(startTime2);
    	var d2 = $.fn.datebox.defaults.parser(value);
    	varify=(d2>=d1);
    	return varify;

    	},
    	message: '结束时间不允许小于开始时间！'
    },
    
    /*验证开始日期小于等于结束日期 */
    mdDate: {
    	validator: function(value, param){
    		startTime2 = '';
		if($(param[0]).length != 0){
			startTime2 = $(param[0]).val();
		};
    	if(!startTime2){
    		var today = new Date();
    		var y = today.getFullYear();  
            var m = today.getMonth()+1;  
            var d = today.getDate();
            var da = (y+"-"+m+"-"+d);
    		startTime2 = da;	
    	}
    	var d1 = $.fn.datebox.defaults.parser(startTime2);
    	var d2 = $.fn.datebox.defaults.parser(value);
    	varify=(d2>=d1);
    	return varify;

    	},
    	message: '结束日期不允许小于开始日期！'
    },
    /*最大指定销售单价不允许小于最小指定销售单价*/
    mdIsSale: {
    	validator: function(value, param){
    		saleNum = '';
		if($(param[0]).length != 0){
			saleNum = $(param[0]).val();
		};
		var maxValue = parseInt(value);
		var minValue =  parseInt(saleNum);
		if(maxValue>=minValue){
			return true;
		}else{
			return false;
		}
    	},
    	message: '最大指定销售单价不允许小于最小指定销售单价！'
    },
    isDuplicateAcount:{
    	validator: function(value, param){
    		return param[0](value);
    	},
    	message:'该帐号已存在！请重新输入'
    },
    /*验证重复，第一个参数为一个方法，通过同步ajax，判断是否重复*/
    isDuplicateName:{
    	validator: function(value, param){
    		return param[0](value);
    	},
    	message:'该名称已存在！请重新输入'
    },
    isDuplicateNumber:{
    	validator: function(value, param){
    		return param[0](value);
    	},
    	message:'该编号已存在！请重新输入'
    },loginName:{
    	validator:function(value,param){
    		 var patrn = /^[a-zA-Z_0-9]+$/;  
        	 if (!patrn.exec(value)) {  
        	 		return false;
        	 }
            return true;
    	},
    	message: '支持字母、数字、下划线'
    },
    floatMoney: {
        validator: function(value){
			// var regu = /(^([1-9]\d*|[0])\.\d{1,2}$|^[1-9]\d*$)/;
			var regu=/^\d+(\.\d{1,2})?$/ ;
			var re = new RegExp(regu);
			if(value > 100000){
				return false;
			}
			return re.test(value);
        },
        message: '请输入[0-100000]的整数且最多保留两位小数'
    },
});