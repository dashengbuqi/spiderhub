/**
 index.js
 Create on 2019/10/25 0025 16:46
 Description: ...
 Create by ZT
 **/
(function () {
    var test_id=0;
    var params ={
        format:$('#format-code'), //格式化代码
        save:$('#save-code'), //保存代码
        font:$('#font-code'), //更换字体的大小
        logErrorType:$('#log-error-type'), //错误类型
        clearResult:$('#clear-result'),//清空结果
        resultArea:$('#result-area'), //结果区
        logArea:$('#log-area'),//日志区
        status:$('#status')
    };
    var opeateParams ={
        run:$('.run'),
        pause:$('.pause'),
        go:$('.go'),
        stop:$('.stop')

    };
    var caches ={
        logs:[],
        result:[],
    };
    var operateBtn = function _opeateBtn() {
        var target = this;
        var _INTERVAL ='';
        var allMessage =[];
        var data ={};
        this.run = function () {
            $.ajax({
                type: 'POST',
                dataType: 'json',
                url: '/clean/begin?id='+app_id,
                data: target._getParams(),
                success: function (res) {
                    if(res.status==1){
                        test_id  = res.data.debug_id;
                        $(params.logArea).append('<p>开始采集数据</p>');
                        caches.logs.push('<p>.........</p>');
                        target.getLogs();
                        $(params.status).removeClass().addClass('status-run').html('正在爬取');
                        $(opeateParams.run).hide();
                        $(opeateParams.proxy).hide();
                        $(opeateParams.pause).show();
                        $(opeateParams.stop).show();
                    }else{
                        art.dialog({title:'提示',time:3, content:res.msg});
                    }
                }
            })
        };
        this.pause = function () {
            clearInterval(_INTERVAL);
            $(params.status).removeClass().addClass('status-pause').html('暂停爬取');
            $(opeateParams.pause).hide();
            $(opeateParams.stop).hide();
            $(opeateParams.run).show();
            $(opeateParams.proxy).show();
            $(opeateParams.go).show();
        };
        this.go = function () {
            this.getLogs();
            $(params.status).removeClass().addClass('status-pause').html('正在爬取');
            $(opeateParams.go).hide();
            $(opeateParams.run).hide();
            $(opeateParams.proxy).hide();
            $(opeateParams.pause).show();
            $(opeateParams.stop).show();

        };
        this.stop = function (){
            $(params.status).removeClass().addClass('status-stop').html('爬虫停止中....');
            $.ajax({
                type: 'put',
                dataType: 'json',
                async:false,
                url: '/clean/end?id='+app_id+'&debug_id='+test_id,
                //data: target._getParams(),
                success: function (data) {
                    if(data.status==1){
                        //  clearInterval(_INTERVAL);
                        $(params.logArea).append('<p>已经点击停止等待系统停止.....</p>');
                        $(params.status).removeClass().addClass('status-stop').html('完成');
                        $(opeateParams.stop).hide();
                        $(opeateParams.run).show();
                        $(opeateParams.pause).hide();
                        $(opeateParams.go).hide();
                    }else{
                        art.dialog({title:'提示',time:3, content:data.msg});
                    }
                }
            })
        };

        this.save =function () {
            $.ajax({
                type: 'post',
                dataType: 'json',
                url: '/clean/save?id='+app_id,
                data: target._getParams(),
                success: function (data) {
                    if(data.status==1){
                        art.dialog({title:'提示',time:3, content:'保存成功'});
                    }else{
                        art.dialog({title:'提示',time:3, content:data.msg});
                    }
                }
            })
        }
        this.renderLogs = function (allData) {
            var str  = '';
            for(var i in allData) {
                allMessage.push(allData[i]);
                if (allData[i].type == 5) {
                    clearInterval(_INTERVAL);
                    $(params.logArea).append('<p>系统已经停止.....</p>');
                    $(params.status).removeClass().addClass('status-stop').html('爬虫已经停止');
                    $(opeateParams.stop).hide();
                    $(opeateParams.run).show();
                    $(opeateParams.proxy).show();
                    $(opeateParams.pause).hide();
                    $(opeateParams.go).hide();
                }
                if (allData[i].content != "") {
                    str += '<p class="logs ' + allData[i].type + '   ' + allData[i].level + '"><i class="circle log-type_' + allData[i].type + '"></i><span class="log-error-type_' + allData[i].level + '">' + allData[i].title + ':' + allData[i].content + '</span></p>';
                }
            }
            if (str != "") {
                $(params.logArea).append(str);
                $(params.logArea)[0].scrollTop =  $(params.logArea)[0].scrollHeight;
            }
        }
        this.renderResult = function (data) {
            var str ='';
            for(var i in data){
                var strItem = JSON.stringify(data[i]);
                str+='<p class="row">【采集结果:】'+strItem+'</p>';
            }
            $(params.resultArea).append(str);
            $(params.resultArea)[0].scrollTop =  $(params.resultArea)[0].scrollHeight;
        }
        this.getLogs = function () {
            _INTERVAL = setInterval(function () {
                $.ajax({
                    type: 'get',
                    dataType: 'json',
                    url: '/clean/heart?id='+app_id+'&debug_id='+test_id,
                    // data: target._getParams(),
                    success: function (res) {
                        var str ='';
                        if(res.data.logs==null&&res.data.rows==null){
                            str ='<span class="waite"></span>';
                        }
                        if(str!=''){
                            $(params.logArea).append(str);
                        }
                        if (res.data.logs != null) {
                            target.renderLogs(res.data.logs);
                        }
                        if (res.data.rows != null) {
                            target.renderResult(res.data.rows);
                        }
                    }
                })
            },2000)
        }
        this._getParams =function () {
            data['code'] = editor.getValue();
            return data;
        }
        this.renderLogByCondition = function () {
            var str ='';
            target._getParams();
            for(var i in allMessage){
                var row =allMessage[i];
                if(data['type']==0&&data['level']==0){
                    str += '<p class="logs ' + row['type'] + ' ' + row['level'] + '"><i class="circle log-type_' + row['type'] + '"></i><span class="log-error-type' + row['level'] + '">' + row['title'] + ':' + row['content'] + '</span></p>';
                }else if(data['type']==0&&data['level']!=0){
                    if(allMessage[i]['level']==data['level']){
                        str += '<p class="logs ' + row['type'] + ' ' + row['level'] + '"><i class="circle log-type_' + row['type'] + '"></i><span class="log-error-type' + row['level'] + '">' + row['title'] + ':' + row['content'] + '</span></p>';
                    }
                }else if(data['log_type']!=0&&data['level']==0){
                    if(allMessage[i]['log_type']==data['type']){
                        str += '<p class="logs ' + row['type'] + ' ' + row['level'] + '"><i class="circle log-type_' + row['type'] + '"></i><span class="log-error-type' + row['level'] + '">' + row['title'] + ':' + row['content'] + '</span></p>';
                    }
                }else{
                    if((allMessage[i]['type']==data['type'])&&allMessage[i]['level']==data['level']){
                        str += '<p class="logs ' + row['type'] + ' ' + row['level'] + '"><i class="circle log-type_' + row['type'] + '"></i><span class="log-error-type' + row['level'] + '">' + row['title'] + ':' + row['content'] + '</span></p>';
                    }
                }
            }
            $(params.logArea).html(str);
            $(params.logArea)[0].scrollTop =  $(params.logArea)[0].scrollHeight;
        }
    }
    var operateFunc = new operateBtn();

    init();
    renderWidth();
    /**
     * 点击事件
     * **/
    function init() {
        $(params.format).click(function () {
            editor.setValue( js_beautify(editor.getValue()));
        });
        $(params.font).change(function () {
            $('.CodeMirror').css('fontSize',$(this).val());
        })
        /**
         * 各种操作*/
        /**
         * 开始爬取的按钮
         */
        $(opeateParams.run).click(function () {
            if(status==1){
                return false;
            }
            operateFunc.run();

        });
        /**
         * 暂停的按钮
         */
        $(opeateParams.pause).click(function () {
            operateFunc.pause();
        });
        /**
         * 继续的按钮
         */
        $(opeateParams.go).click(function () {
            operateFunc.go();
        });
        /**
         * 停止的按钮
         */
        $(opeateParams.stop).click(function () {
            operateFunc.stop();
        });
        /**
         * 清空
         */
        $(params.clearResult).click(function () {
            $(params.resultArea).html('');
        });
        /**
         * 保存
         */
        $(params.save).click(function () {
            operateFunc.save();
        });

        /**
         * 切换
         */
        $(params.logType).change(function () {
            var value  = $(this).val();
            operateFunc.renderLogByCondition();

        });
        $(params.logErrorType).change(function () {
            var value  = $(this).val();
            operateFunc.renderLogByCondition();
        });
    }


    /**
     * 自动改变左侧的宽度
     */
    function renderWidth() {
        $('#mainSplitter').jqxSplitter({ width: '100%', height: '100%', panels: [{ size: '70%' }, { size: '30%'}] });
        $('#rightSplitter').jqxSplitter({ height: '100%', orientation: 'horizontal', panels: [{ size: '50%', collapsible: false }, { size: '50%'}] });
    }
})();