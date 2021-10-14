(function(E, C, D, A) {
    var B, $, _, J = "@ARTDIALOG.DATA", K = "@ARTDIALOG.OPEN", H = "@ARTDIALOG.OPENER", I = C.name = C.name || "@ARTDIALOG.WINNAME" + (new Date).getTime(), F = C.VBArray && !C.XMLHttpRequest;
    E(function() {
        !C.jQuery && document.compatMode === "BackCompat" && alert("artDialog Error: document.compatMode === \"BackCompat\"")
    });
    var G = D.top = function() {
        var _ = C
          , $ = function(A) {
            try {
                var _ = C[A].document;
                _.getElementsByTagName
            } catch ($) {
                return !1
            }
            return C[A].artDialog && _.getElementsByTagName("frameset").length === 0
        }
        ;
        return $("top") ? _ = C.top : $("parent") && (_ = C.parent),
        _
    }();
    D.parent = G,
    B = G.artDialog,
	B.defaults.zIndex = 10000,
    _ = function() {
        return B.defaults.zIndex
    }
    ,
    D.data = function(C, B) {
        var _ = D.top
          , $ = _[J] || {};
        _[J] = $;
        if (B !== A)
            $[C] = B;
        else
            return $[C];
        return $
    }
    ,
    D.removeData = function(_) {
        var $ = D.top[J];
        $ && $[_] && delete $[_]
    }
    ,
    D.through = $ = function() {
        var $ = B.apply(this, arguments);
        return G !== C && (D.list[$.config.id] = $),
        $
    }
    ,
    G !== C && E(C).bind("unload", function() {
        var A = D.list, _;
        for (var $ in A)
            A[$] && (_ = A[$].config,
            _ && (_.duration = 0),
            A[$].close(),
            delete A[$])
    }),
    D.open = function(B, P, O) {
        P = P || {};
        var N, L, M, X, W, V, U, T, S, R = D.top, Q = "position:absolute;left:-9999em;top:-9999em;border:none 0;background:transparent", a = "width:100%;height:100%;border:none 0";
        if (O === !1) {
            var Z = (new Date).getTime()
              , Y = B.replace(/([?&])_=[^&]*/, "$1_=" + Z);
            B = Y + (Y === B ? (/\?/.test(B) ? "&" : "?") + "_=" + Z : "")
        }
        var G = function() {
            var B, C, _ = L.content.find(".aui_loading"), A = N.config;
            M.addClass("aui_state_full"),
            _ && _.hide();
            try {
                T = W.contentWindow,
                U = E(T.document),
                S = T.document.body
            } catch ($) {
                W.style.cssText = a,
                A.follow ? N.follow(A.follow) : N.position(A.left, A.top),
                P.init && P.init.call(N, T, R),
                P.init = null ;
                return
            }
            B = A.width === "auto" ? U.width() + (F ? 0 : parseInt(E(S).css("marginLeft"))) : A.width,
            C = A.height === "auto" ? U.height() : A.height,
            setTimeout(function() {
                W.style.cssText = a
            }, 0),
            N.size(B, C),
            A.follow ? N.follow(A.follow) : N.position(A.left, A.top),
            P.init && P.init.call(N, T, R),
            P.init = null
        }
          , I = {
            zIndex: _(),
            init: function() {
                N = this,
                L = N.DOM,
                X = L.main,
                M = L.content,
                W = N.iframe = R.document.createElement("iframe"),
                W.src = B,
                W.name = "Open" + N.config.id,
                W.style.cssText = Q,
                W.setAttribute("frameborder", 0, 0),
                W.setAttribute("allowTransparency", !0),
                V = E(W),
                N.content().appendChild(W),
                T = W.contentWindow;
                try {
                    T.name = W.name,
                    D.data(W.name + K, N),
                    D.data(W.name + H, C)
                } catch ($) {}
                V.bind("load", G)
            },
            close: function() {
                V.css("display", "none").unbind("load", G);
                if (P.close && P.close.call(this, W.contentWindow, R) === !1)
                    return !1;
                M.removeClass("aui_state_full"),
                V[0].src = "about:blank",
                V.remove();
                try {
                    D.removeData(W.name + K),
                    D.removeData(W.name + H)
                } catch ($) {}
            }
        };
        typeof P.ok == "function" && (I.ok = function() {
            return P.ok.call(N, W.contentWindow, R)
        }
        ),
        typeof P.cancel == "function" && (I.cancel = function() {
            return P.cancel.call(N, W.contentWindow, R)
        }
        ),
        delete P.content;
        for (var J in P)
            I[J] === A && (I[J] = P[J]);
        return $(I)
    }
    ,
    D.open.api = D.data(I + K),
    D.opener = D.data(I + H) || C,
    D.open.origin = D.opener,
    D.close = function() {
        var $ = D.data(I + K);
        return $ && $.close(),
        !1
    }
    ,
    G != C && E(document).bind("mousedown", function() {
        var $ = D.open.api;
        $ && $.zIndex()
    }),
    D.load = function(C, D, B) {
        B = B || !1;
        var G = D || {}
          , H = {
            zIndex: _(),
            init: function(A) {
                var _ = this
                  , $ = _.config;
                E.ajax({
                    url: C,
                    success: function($) {
                        _.content($),
                        G.init && G.init.call(_, A)
                    },
                    cache: B
                })
            }
        };
        delete D.content;
        for (var F in G)
            H[F] === A && (H[F] = G[F]);
        return $(H)
    }
    ,
    D.alert = function(B, A) {
        return $({
            id: "Alert",
            zIndex: _(),
            icon: "warning",
            fixed: !0,
            lock: !0,
            content: B,
            ok: !0,
            close: A
        })
    }
    ,
    D.confirm = function(C, A, B) {
        return $({
            id: "Confirm",
            zIndex: _(),
            icon: "question",
            fixed: !0,
            lock: !0,
            opacity: 0.1,
            content: C,
            ok: function($) {
                return A.call(this, $)
            },
            cancel: function($) {
                return B && B.call(this, $)
            }
        })
    }
    ,
    D.prompt = function(D, B, C) {
        C = C || "";
        var A;
        return $({
            id: "Prompt",
            zIndex: _(),
            icon: "question",
            fixed: !0,
            lock: !0,
            opacity: 0.1,
            content: ["<div style=\"margin-bottom:5px;font-size:12px\">", D, "</div>", "<div>", "<input value=\"", C, "\" style=\"width:18em;padding:6px 4px\" />", "</div>"].join(""),
            init: function() {
                A = this.DOM.content.find("input")[0],
                A.select(),
                A.focus()
            },
            ok: function($) {
                return B && B.call(this, A.value, $)
            },
            cancel: !0
        })
    }
    ,
    D.tips = function(B, A) {
        return $({
            id: "Tips",
            zIndex: _(),
            title: !1,
            cancel: !1,
            fixed: !0,
            lock: !1
        }).content("<div style=\"padding: 0 1em;\">" + B + "</div>").time(A || 1.5)
    }
    ,
    E(function() {
        var A = D.dragEvent;
        if (!A)
            return;
        var B = E(C)
          , $ = E(document)
          , _ = F ? "absolute" : "fixed"
          , H = A.prototype
          , I = document.createElement("div")
          , G = I.style;
        G.cssText = "display:none;position:" + _ + ";left:0;top:0;width:100%;height:100%;" + "cursor:move;filter:alpha(opacity=0);opacity:0;background:#FFF",
        document.body.appendChild(I),
        H._start = H.start,
        H._end = H.end,
        H.start = function() {
            var E = D.focus.DOM
              , C = E.main[0]
              , A = E.content[0].getElementsByTagName("iframe")[0];
            H._start.apply(this, arguments),
            G.display = "block",
            G.zIndex = D.defaults.zIndex + 3,
            _ === "absolute" && (G.width = B.width() + "px",
            G.height = B.height() + "px",
            G.left = $.scrollLeft() + "px",
            G.top = $.scrollTop() + "px"),
            A && C.offsetWidth * C.offsetHeight > 307200 && (C.style.visibility = "hidden")
        }
        ,
        H.end = function() {
            var $ = D.focus;
            H._end.apply(this, arguments),
            G.display = "none",
            $ && ($.DOM.main[0].style.visibility = "visible")
        }
    })
})(this.art || this.jQuery, this, this.artDialog)
