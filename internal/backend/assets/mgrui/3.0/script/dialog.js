/**
 * ajax确认对话框
 * @param title 对话框标题
 * @param msg 对话框信息
 * @param fn 回调函数，函数的形式为 fn(b)，b为boolean类型，指明用户选择了确认(true)还是取消(false)
 * @return void
 */
function msgConfirm(title, msg, fn){
  if (fn == null){
    fn = nullFnC;
  }
  $.messager.confirm(title,msg,fn);
}

/**
 * ajax信息框
 * @param title 对话框标题
 * @param msg 对话框信息
 * @param icon 图标，'error','question','info','warning'
 * @param fn 回调函数
 * @return
 */
function msgAlert(title, msg, icon, fn){
  if (fn == null){
    fn = nullFn;
  }
  $.messager.alert(title, msg, icon, fn);
}

/**
 * ajax输入对话框
 * @param title 对话框标题
 * @param msg 对话框信息
 * @param fn 回调函数，函数的形式为 fn(val)，val为用户在对话框中输入的内容
 * @return void
 */
function msgInput(title, msg, fn){
  if (fn == null){
    fn = nullFnC;
  }
  $.messager.prompt(title,msg,fn);
}

function nullFn(){  
}

function nullFnC(b){
  return b;
}