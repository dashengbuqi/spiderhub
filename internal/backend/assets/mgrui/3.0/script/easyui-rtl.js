(function($){
	$.fn.datagrid.defaults.resizeHandle = 'left';
	
	var datagrid_freezeRow = $.fn.datagrid.methods.freezeRow;
	$.fn.datagrid.methods.freezeRow = function(jq, index){
		return jq.each(function(){
			datagrid_freezeRow.call(this, jq, index);
			var state = $.data(this, 'datagrid');
			if (!state.rtlscroll){
				state.rtlscroll = true;
				var dc = state.dc;
				dc.body2.bind('scroll', function(){
					var ftable = $(this).children('table.datagrid-btable-frozen');
					ftable.css('left',  $(this)._outerWidth() + $(this)._scrollLeft() - ftable._outerWidth());
				});
			}
		});
	}
	

	$.fn._scrollLeft = function(left){
		if (left == undefined){
			if ($.browser.msie){
				return this.scrollLeft();
			} else if ($.browser.mozilla){
				return -this.scrollLeft();
			} else {
				return this[0].scrollWidth - this[0].clientWidth - this.scrollLeft();
			}
		}
		return this.each(function(){
			if ($.browser.msie){
				$(this).scrollLeft(left);
			} else if ($.browser.mozilla){
				$(this).scrollLeft(-left);
			} else {
				$(this).scrollLeft(this.scrollWidth - this.clientWidth - left);
			}
		});
	}
	
	
})(jQuery);
var userAgent = navigator.userAgent.toLowerCase();
//Figure out what browser is being used
jQuery.browser = {
version: (userAgent.match( /.+(?:rv|it|ra|ie)[\/: ]([\d.]+)/ ) || [])[1],
safari: /webkit/.test( userAgent ),
opera: /opera/.test( userAgent ),
msie: /msie/.test( userAgent ) && !/opera/.test( userAgent ),
mozilla: /mozilla/.test( userAgent ) && !/(compatible|webkit)/.test( userAgent )
}; 