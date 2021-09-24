// JavaScript Document
function itemsShortcut(option) {

	var self = this;
	option = $.extend({
		container : $(".itemsShortcut"),
		items : $(".loadContent .productItem"),
		speed : 100,
		scrollMain : $('.newProductContent').parent()
	}, option || {});

	/* 初始化参数 */
	self.init = function() {
		self.container = option.container, self.items = option.items,
				self.speed = option.speed,
				self.itemsLength = option.items.length,
				self.scrollMain = option.scrollMain;
	};
	self.init();
	/* 创建栏目快捷主体 */
	self.createdShortcut = function() {
		self.container.css({
			"position" : "fixed",
			top : "50%",
			right : "17px",
			"z-index" : "1000"
		}).empty();
		$(
				'<div class="shortcutTarget"><div class="shortcutIcon">&lt;&lt;</div>快速定位</div>')
				.appendTo(self.container);
		$('<ul></ul>').appendTo(self.container);
		for ( var i = 0; i < self.itemsLength; i++) {
			var cutList = $(
					'<li><div class="shortcutItem"><div class="shortcutNum">'
							+ (i + 1) + '</div><div class="shortcutTitle">'
							+ self.items.eq(i).find(".itemTitle11 >span").text()
							+ '</div></div></li>').appendTo(
					self.container.find("ul"));
		}
		;
		self.container.css({
			"margin-top" : -self.container.height() / 2 + "px"
		})

	};
	self.createdShortcut();
	/* 点击事件 */
	self.container.find("li").bind("click",function(e) {
				e.preventDefault();
				var scTop = self.items.eq($(this).index()).position().top;
				self.items.eq($(this).index()).find(".itemContent").slideDown();
				self.items.eq($(this).index()).find(".itemTitle").removeClass("active");
				self.scrollMain.animate({
					scrollTop : scTop
				}, self.speed);
			})
}