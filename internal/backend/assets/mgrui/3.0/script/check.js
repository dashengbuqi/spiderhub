
 /**
     *检查电话号码格式
     *如果格式不正确     返回false  正确   返回true
     */
    function checkPhoneFormat(object){
      var pattern = /^[0-9-]{6,40}$/;
       flag = pattern.test(object.value);
       if(object.value=="")
         return true;
        if(!flag) {
        alert('联系电话格式不正确!');
		object.focus();
		return false;
       } else {
        return true;
       }
     }
     /**
     *检查电话号码格式
     *如果格式不正确     返回false  正确   返回true
     */

 /*(1)电话号码由数字、"("、")"和"-"构成
　　(2)电话号码为3到8位
　　(3)如果电话号码中包含有区号，那么区号为三位或四位
　　(4)区号用"("、")"或"-"和其他部分隔开
　　(5)移动电话号码为11或12位，如果为12位,那么第一位为0
　　(6)11位移动电话号码的第一位和第二位为"13"
　　(7)12位移动电话号码的第二位和第三位为"13"
     */
     function PhoneCheck(object) {
       var str=object.value;
       var reg=/(^[0-9]{3,4}\-[0-9]{3,8}$)|(^[0-9]{3,8}$)|(^\([0-9]{3,4}\)[0-9]{3,8}$)|(^0{0,1}13[0-9]{9}$)/;
       flag = reg.test(object.value);
       if(object.value=="")
         return true;
       if(!flag) {
        alert('电话格式不正确!');
		object.focus();
		return false;
       } else {
        return true;
       }
     }
    //检查是否数字
	function checkNUM(NUM)
	{
	 var i,j,strTemp;
	 strTemp="0123456789";
	 if ( NUM.length== 0)
	  return 0
	 for (i=0;i<NUM.length;i++)
	 {
	  j=strTemp.indexOf(NUM.charAt(i));
	  if (j==-1)
	  {
	   return 0;
	  }
	 }
	 return 1;
	}

	//检查是否年（YYYY）
	function chkYear(datestr)
	{
	 var lthdatestr
	 if (datestr != "")
	  lthdatestr= datestr.length ;
	 else
	  return 1;

	 if (lthdatestr== 0)
	   return 1;
	 if (lthdatestr != 4)
	   return 0;

	 if (checkNUM(datestr) == 0)
	   return 0;

	 return 1;
	}

	//检查是否月份（YYYYMM）
	function chkMonth(datestr)
	{
	 var lthdatestr
	 if (datestr != "")
	  lthdatestr= datestr.length ;
	 else
	  return 1;

	 if (lthdatestr== 0)
	   return 1;
	 if (lthdatestr != 6)
	   return 0;

	 if (checkNUM(datestr) == 0)
	   return 0;

	 tmp = new String(datestr);

	 var month= tmp.substr(4,2);
	 if (!((1<=month) && (12>=month)) )
	 {
	    return 0;
	 }
	 return 1;
	}


	//检查是否日期（YYYYMMDD）
	function chkdate(datestr)
	{
	 var lthdatestr
	 if (datestr != "")
	  lthdatestr= datestr.length ;
	 else
	  return 1;

	 if (lthdatestr== 0)
	   return 1;
	 if (lthdatestr != 8)
	   return 0;

	 if (checkNUM(datestr) == 0)
	   return 0;

	 tmp = new String(datestr);

	 var year = tmp.substr(0,4);
	 var month= tmp.substr(4,2);
	 var day= tmp.substr(6,2);

	 if (!((1<=month) && (12>=month) && (31>=day) && (1<=day)) )
	 {
	  return 0;
	 }
	 if (!((year % 4)==0) && (month==2) && (day==29))
	 {
	  return 0;
	 }
	 if ((month<=7) && ((month % 2)==0) && (day>=31))
	 {
	  return 0;

	 }
	 if ((month>=8) && ((month % 2)==1) && (day>=31))
	 {
	  return 0;
	 }
	 if ((month==2) && (day==30))
	 {
	  return 0;
	 }

	 return 1;
	}


    //----去空格----
    function LTrim(str){
    	var i=0;
    	while ( str.charAt(i) == ' ') i++;
    	return str.substring(i);
    }
    function RTrim(str){
    	var i=0;
    	var len = str.length;
    	while (str.charAt(len-i-1) == ' ') i++;
    	return str.substring(0,len-i);
    }
    function ATrim(str){
    	return RTrim(LTrim(str));
    }

    //----判断数据项非空----
    function isEmpty(str){
      if(str != null){
        var tmp = ATrim(str)
        if(tmp.length == 0){
            return true;
        }
      }
      return false;
    }

    //----判断数据项输入长度是否不够----
    function isLenErr(str, len){
        var tmp = ATrim(str)
    	if(tmp.length == len){
            return false;
        }
        return true;
    }

    //----判断数据项输入长度是否超长----
    function isLenOver(str, len){
        var tmp = ATrim(str)
    	if(tmp.length <= len){
            return false;
        }
        return true;
    }

    //----检查费率是否错误----
    function isRateErr(str){
      var tmp = ATrim(str)
      var bOneDot = false;
      for(var i = tmp.length-1; i >= 0 ; i--){
          var oneChar = tmp.charAt(i)
          if(oneChar == "." && (!bOneDot)){
              if((tmp.length-1-i)>4)
                  return true;

              bOneDot = true;
              continue;
          }
          if(oneChar < "0" || oneChar > "9"){
      	    return true;
          }
      }
      if(tmp > 1){
        return true;
      }
      return false;
    }

    //----检查金额是否错误----
    function isAmountErr(str){
      var tmp = ATrim(str)
      var bOneDot = false;
      for(var i = tmp.length-1; i >= 0 ; i--){
        var oneChar = tmp.charAt(i)
        if(oneChar == "." && (!bOneDot)){
          if((tmp.length-1-i)>2){
            return true;
          }

          bOneDot = true;
          continue;
        }
        if(oneChar < "0" || oneChar > "9"){
    	    return true;
        }
      }
      return false;
    }

    //----判断是否闰年----
    function isRyear(inputInt){
        if (inputInt % 100 == 0 && inputInt % 400 == 0 || inputInt % 100 != 0 && inputInt % 4 == 0){
            return true
        }else{
	        return false
        }
    }

    //----判断日期是否合法----
    function isDate(str){
        var year = parseFloat(str.substring(0,4))
        var month = parseFloat(str.substring(4,6))
        var day = parseFloat(str.substring(6,8))
        if (month < 1 || month > 12 || day < 1 || day > 31 || year < 1000 || year > 2050){
            return false
        }else if ((month == 4 || month == 6 || month == 9 || month ==11) && (day > 30)){
	        return false
	    }else if (isRyear(year) && month == 2 && day > 29 || !isRyear(year) && month == 2 && day > 28){
	        return false
	    }else{
	        return true
	    }
    }

    //----判断输入项类型是否合法日期----
    function isCorrectDate(str){
        if (str.length != 8) {
            return false
    	}
        for(var i = 0;i < 8;i++){
            var oneChar = str.charAt(i)
            if(oneChar < "0" || oneChar > "9"){
                return false
            }
    	}
        if (!isDate(str)){
            return false
        }else{
            return true
    	}
    }
    //--------------检查是否合法手机号码------------------
		function isMobile(mobile){

			if ( mobile.length != 11 ){
				return false;
			}
			if( checkNUM(mobile) == 0){
				return false;
			}
	 		var s = mobile.substring(0,2);
			if( s != "13" && s != "15" && s != "18"){
				return false;
			}
			return true;

	 }
// -->
    //----判断输入项是否由数字或字母或数字和字母组成----
	function isNumOrLetter(str){
		var regu = "^[0-9a-zA-Z]+$";

		var re = new RegExp(regu);

		if (re.test(str)) {

			return true;

		} else {

			return false;

		}
	}
	
	
	//正整数和最多包含两位小数点的正数，不包含0
	function isAmount(amount){
		if(amount == '0') {
			return 1;
		}
		var regu = /(^([1-9]\d*|[0])\.\d{1,2}$|^[1-9]\d*$)/;
		var re = new RegExp(regu);
		return re.test(amount);
	}
	
	//true  : 非法返回true
	//false : 正常返回false
	//--------------检查邮箱号码是否合法------------------
	function Email(email){
		var emailReg2 = /^.+\@(\[?)[a-zA-Z0-9\-\.]+\.([a-zA-Z]{2,6}|[0-9]{1,3})(\]?)$/;
		if (emailReg2.test(email)){
			return true;
		}
		else{
			return false;
		}
		
	}
	

          // 身份证验证位值.10代表X 

    function Isidcard(num) { 
    	 num = num.toUpperCase();           //身份证号码为15位或者18位，15位时全为数字，18位前17位为数字，最后一位是校验位，可能为数字或字符X。        
    	    if (!(/(^\d{15}$)|(^\d{17}([0-9]|X)$)/.test(num))) {     
    	       // alert('输入的身份证号长度不对，或者号码不符合规定！\n15位号码应全为数字，18位号码末位可以为数字或X。');              
    	        return false;         
    	    } //校验位按照ISO 7064:1983.MOD 11-2的规定生成，X可以认为是数字10。 
    	    //下面分别分析出生日期和校验位 
    	    var len, re; len = num.length; if (len == 15) { 
    	        re = new RegExp(/^(\d{6})(\d{2})(\d{2})(\d{2})(\d{3})$/); 
    	        var arrSplit = num.match(re);  //检查生日日期是否正确
    	        var dtmBirth = new Date('19' + arrSplit[2] + '/' + arrSplit[3] + '/' + arrSplit[4]); 
    	        var bGoodDay; bGoodDay = (dtmBirth.getYear() == Number(arrSplit[2])) && ((dtmBirth.getMonth() + 1) == Number(arrSplit[3])) && (dtmBirth.getDate() == Number(arrSplit[4]));
    	        if (!bGoodDay) {         
    	           // alert('输入的身份证号里出生日期不对！');            
    	            return false;
    	        } else { //将15位身份证转成18位 //校验位按照ISO 7064:1983.MOD 11-2的规定生成，X可以认为是数字10。        
    	            var arrInt = new Array(7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2);         
    	            var arrCh = new Array('1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2');      
    	            var nTemp = 0, i;            
    	            num = num.substr(0, 6) + '19' + num.substr(6, num.length - 6);           
    	            for(i = 0; i < 17; i ++) {                 
    	                nTemp += num.substr(i, 1) * arrInt[i];        
    	            }
    	            num += arrCh[nTemp % 11]; 
    	            return true;
    	        }
    	    }
    	    if (len == 18) {
    	        re = new RegExp(/^(\d{6})(\d{4})(\d{2})(\d{2})(\d{3})([0-9]|X)$/); 
    	        var arrSplit = num.match(re);  //检查生日日期是否正确 
    	        var dtmBirth = new Date(arrSplit[2] + "/" + arrSplit[3] + "/" + arrSplit[4]); 
    	        var bGoodDay; bGoodDay = (dtmBirth.getFullYear() == Number(arrSplit[2])) && ((dtmBirth.getMonth() + 1) == Number(arrSplit[3])) && (dtmBirth.getDate() == Number(arrSplit[4])); 
    	        if (!bGoodDay) { 
    	            //alert(dtmBirth.getYear()); 
    	           // alert(arrSplit[2]); 
    	           /// alert('输入的身份证号里出生日期不对！'); 
    	            return false; 
    	        }
    	        else { //检验18位身份证的校验码是否正确。 //校验位按照ISO 7064:1983.MOD 11-2的规定生成，X可以认为是数字10。 
    	            var valnum; 
    	            var arrInt = new Array(7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2); 
    	            var arrCh = new Array('1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'); 
    	            var nTemp = 0, i; 
    	            for(i = 0; i < 17; i ++) { 
    	                nTemp += num.substr(i, 1) * arrInt[i];
    	            } 
    	            valnum = arrCh[nTemp % 11]; 
    	            if (valnum != num.substr(17, 1)) { 
    	               // alert('18位身份证的校验码不正确！应该为：' + valnum); 
    	                return false; 
    	            } 
    	            return true; 
    	        } 
    	    } return false;  
    } 

