/// <reference path="jquery-1.7.min.js" />

/*
    进行图片预览功能的JS
*/
(function () {
    var regURL = /(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?/i;

    var enablePreview = $.browser.mozilla || $.browser.msie;
    
    $.fn.extend({
        preview: function (opts) {
            /// <summary>设置图片的预览，忽略非file元素</summary>
            /// <param name="opts" type="object">预览时所用的参数，包含 viewContain属性： 一个img元素集合，用来显示预览图</param>
            
            //如果参数是字符串，则认为该字符串是图片地址，直接进行预览
            if (typeof (opts) == "string") {
           
                //只有 img 元素才进行预览
                if (this.length >= 0 && this[0].nodeName == "IMG") {
                    //如果图片地址是一个URL或者是在火狐浏览器中，则直接设置img元素的src属性
                    if ($.browser.mozilla || regURL.test(opts)) {
                        this.show().attr("src",opts);
                    }else if($.browser.msie){
                        try {
                            //IE中使用滤镜，在 img 的父元素上显示预览
                            var viewCtrl = this.parent()[0];
                            
                            viewCtrl.style.filter = "progid:DXImageTransform.Microsoft.AlphaImageLoader(sizingMethod=scale)";
                            viewCtrl.filters.item("DXImageTransform.Microsoft.AlphaImageLoader").src = opts;
                            this.hide();
                        } catch (e) {
                            //alert("您上传的图片格式不正确，请重新选择!");
                            //return false;
                        }
                    }
                    else{
                    this.show().attr("src",opts);
                    }
                }
                
                return this;
            }
            //缺少参数或者不是IE、Firefox 浏览器，则没有预览功能
            if (!opts.viewContain) {
            
                return this;
            }

            //获取用来显示预览的元素数量
            var viewContainerCount = opts.viewContain.length;

            //根据 viewContainerCount 设置file中可进行预览的元素
            var files = this.filter(":file:lt(" + viewContainerCount + ")").change(function () {
            
                //当前file元素在集合中的哪个位置
                var index = files.index(this);
                
                //获取file元素所在位置对应的预览元素
                var viewCtrl = opts.viewContain[index];
                
                if (!viewCtrl) {
                    return;
                }


                //var width = opts.width;
                //var height = opts.height;
                
                //检查是否支持html中的属性
                if (this.files && this.files[0]) {
                    //火狐下，直接设img属性
                    //imgObjPreview.style.display = 'block';
                    //imgObjPreview.style.width = '300px';
                    //imgObjPreview.style.height = '120px';
                    //viewCtrl.src = this.files[0].getAsDataURL();
                     
                    //火狐7以上版本不能用上面的getAsDataURL()方式获取，需要一下方式  
                    if(window.URL){
                    viewCtrl.src = window.URL.createObjectURL(this.files[0]);
                    }else{
                    //周祥扩展 google 360预览 
                     viewCtrl.src = window.webkitURL.createObjectURL(this.files[0]);
                    }
                   
                    $(viewCtrl).attr("url",viewCtrl.src);
                }
                else {
                    //隐藏 img元素，避免 img.src 的图片覆盖滤镜
                    viewCtrl.style.display = "none";

                    

                    //选择图片的绝对路径
                    this.select();

                    //在IE9下，如果file控件获得焦点，则document.selection.createRange()拒绝访问
                    this.blur();

                    //获取绝对路径
                    var imgSrc = document.selection.createRange().text;

                    //var localImagId = document.getElementById("localImag");
                    //必须设置初始大小
                    //localImagId.style.width = "300px";
                    //localImagId.style.height = "120px";
                    //图片异常的捕捉，防止用户修改后缀来伪造图片
                    try {
                        
                        $(viewCtrl).attr("url",imgSrc);
                        
                        //将显示预览的元素设置为 img的父元素
                        viewCtrl = viewCtrl.parentNode;

                        //使用滤镜
                        viewCtrl.style.filter = "progid:DXImageTransform.Microsoft.AlphaImageLoader(sizingMethod=scale)";
                        viewCtrl.filters.item("DXImageTransform.Microsoft.AlphaImageLoader").src = imgSrc;
                    } catch (e) {
                        //alert("您上传的图片格式不正确，请重新选择!");
                        //return false;
                    }
                    //imgObjPreview.style.display = 'none';
                    document.selection.empty();
                }
            });

            return this;
        },
        clearPreview: function(opts){
            if (!opts) {
                opts = {};
            }
            if (typeof(opts) == "string") {
                opts = { defautUrl: opts };
            }

            if (this.length >= 0 && this[0].nodeName == "IMG") {
                this.attr("src","").removeAttr("url");
                if($.browser.msie){
                    this.parent().css("filter","");
                }
                opts.defautUrl && this.preview(opts.defautUrl);
            }
            return this;
        }
    });
})();