/* 
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
var GlobalArrTab= [];//全局tab数组
var GlobalTabsId="tabs";//全局tabs标识id
//表单封装
$.fn.custom_form= function(opt_params)
{
    if($(this).form('validate')){
        showProcess(true, '温馨提示', '正在提交数据...');
        $.ajax({
            url : opt_params.url,
            data:$(this).serialize(),
            type : "POST",
            dataType:'json',
            success : function(data) {
                var json = data;//eval('('+data+')');
                showProcess(false);
                 if (json.status) {
                     if (opt_params.is_handle == 'close') {
                         var title = parent.$('#tabs').tabs('getSelected').panel('options').title;
                         var parent_title = window.parent.$("iframe[name=" + title + "]").attr('parent_name');
                         if (parent_title) {
                             var tabs = getTabObj();
                             if(opt_params.refresh){  //刷新当前页
                                 refreshTab(parent_title);
                             }
                             tabs.tabs('select', parent_title);
                         }
                         var callback = parent.closeTab(title);
                         if (jQuery.isFunction(callback)) {
                             callback.call();
                         }
                     }else if (opt_params.is_handle === 'back' && !opt_params.back_url){
                         history.back(-1);return false;
                     } else if (opt_params.back_url) {
                         window.location = opt_params.back_url; return false;
                     } else{
                         showMsg(json.msg,true);
                     }
                 }else {
                    showMsg(json.msg, false, false, "error");
                 }
            }
        });
    }
};

/**
 * 重写后 解决调用时 右上角关闭无效问题
 */
function custom_confirm_ajax(title,msg,url,ajax_type,reload){
    if(ajax_type==null){
        ajax_type="post";
    }
    if(reload == null){
        reload = true;
    }

   var messagebody= showConfirm(title,msg,function(){        
        showProcess(true, '温馨提示', '正在提交数据...');
        $.ajax({
            url:url,
            type:ajax_type,
            dataType:"json",
            success : function(data) {
                var json = data;
                showProcess(false);
                if (reload){
                    reloadDatagrid();
                }
                var iconStr = "info";
                if (json.status == 0) {
                    iconStr =  "error"
                }
                $.messager.alert('温馨提示',json.msg,iconStr);

            }
        });
    });

    var messagebox=messagebody.parent();
    var closeBox =messagebox.find(".panel-tool-close");
        closeBox.unbind("click").bind("click",function(){
            //closeMyWindow();
          //  reloadDatagrid();
        //  location.reload();
            messagebox.siblings(".window-shadow").remove();
            messagebox.siblings(".window-mask").remove();
            messagebox.remove();
        });  
}
jQuery.fn.extend({
    everyTime: function(interval, label, fn, times, belay) {
        return this.each(function() {
            jQuery.timer.add(this, interval, label, fn, times, belay);
        });
    },
    oneTime: function(interval, label, fn) {
        return this.each(function() {
            jQuery.timer.add(this, interval, label, fn, 1);
        });
    },
    stopTime: function(label, fn) {
        return this.each(function() {
            jQuery.timer.remove(this, label, fn);
        });
    }
});


jQuery.extend({
    timer: {
        guid: 1,
        global: {},
        regex: /^([0-9]+)\s*(.*s)?$/,
        powers: {
            // Yeah this is major overkill...
            'ms': 1,
            'cs': 10,
            'ds': 100,
            's': 1000,
            'das': 10000,
            'hs': 100000,
            'ks': 1000000
        },
        timeParse: function(value) {
            if (value == undefined || value == null)
                return null;
            var result = this.regex.exec(jQuery.trim(value.toString()));
            if (result[2]) {
                var num = parseInt(result[1], 10);
                var mult = this.powers[result[2]] || 1;
                return num * mult;
            } else {
                return value;
            }
        },
        add: function(element, interval, label, fn, times, belay) {
            var counter = 0;

            if (jQuery.isFunction(label)) {
                if (!times)
                    times = fn;
                fn = label;
                label = interval;
            }

            interval = jQuery.timer.timeParse(interval);

            if (typeof interval != 'number' || isNaN(interval) || interval <= 0)
                return;

            if (times && times.constructor != Number) {
                belay = !!times;
                times = 0;
            }

            times = times || 0;
            belay = belay || false;

            if (!element.$timers)
                element.$timers = {};

            if (!element.$timers[label])
                element.$timers[label] = {};

            fn.$timerID = fn.$timerID || this.guid++;

            var handler = function() {
                if (belay && this.inProgress)
                    return;
                this.inProgress = true;
                if ((++counter > times && times !== 0) || fn.call(element, counter) === false)
                    jQuery.timer.remove(element, label, fn);
                this.inProgress = false;
            };

            handler.$timerID = fn.$timerID;

            if (!element.$timers[label][fn.$timerID])
                element.$timers[label][fn.$timerID] = window.setInterval(handler,interval);

            if ( !this.global[label] )
                this.global[label] = [];
            this.global[label].push( element );

        },
        remove: function(element, label, fn) {
            var timers = element.$timers, ret;

            if ( timers ) {

                if (!label) {
                    for ( label in timers )
                        this.remove(element, label, fn);
                } else if ( timers[label] ) {
                    if ( fn ) {
                        if ( fn.$timerID ) {
                            window.clearInterval(timers[label][fn.$timerID]);
                            delete timers[label][fn.$timerID];
                        }
                    } else {
                        for ( var fn in timers[label] ) {
                            window.clearInterval(timers[label][fn]);
                            delete timers[label][fn];
                        }
                    }

                    for ( ret in timers[label] ) break;
                    if ( !ret ) {
                        ret = null;
                        delete timers[label];
                    }
                }

                for ( ret in timers ) break;
                if ( !ret )
                    element.$timers = null;
            }
        }
    }
});
