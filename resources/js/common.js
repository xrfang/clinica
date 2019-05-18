function toast(title, mesg, icon) {
    if (typeof (icon) == 'undefined') {
        icon = 'warning'
    }
    $.toast({
        heading: '<b style="font-size:15px;font-family:courier">' + title + "</b>",
        text: '<span style="font-size:15px;font-family:courier">' + mesg + '</span>',
        position: 'mid-center',
        icon: icon,
        loader: false,
        bgColor: '#666666',
        textColor: 'white',
        textAlign: 'left',
        allowToastClose: false,
        stack: false,
        hideAfter: 5000,
        showHideTransition: 'fade'
    })
}

$(document).ajaxError(function (_, xhr, settings) {
    if (xhr.status >= 200 && xhr.status < 300) {
        return //2xx reply means everything's ok
    }
    // var stat = xhr.status + " " + xhr.statusText
    // var mesg = settings.type + ' ' + settings.url + '<br><p style="white-space:pre">' + xhr.responseText + '</p>'
    var stat = `操作失败`
    var mesg = `<p>${xhr.responseText}</p>`
    toast(stat, mesg)
});

function caseBgColor(status) {
    switch (parseInt(status)) {
        case 0: //尚未结束
            return "primary"
        case 1: //痊愈/显效
            return "success"
        case 2: //失败
            return "danger"
        case 3: //无反馈
            return "secondary"
    }
}

function consultBgColor(status) {
    switch (parseInt(status)) {
        case 0: //就诊完成
            return "#d4edda" //success
        case 1: //预约中
            return "#cce5ff" //primary
        case 2: //未赴约
            return "#f8d7da" //danger
        case 3: //取消
            return "#e2e3e5" //secondary
    }
}

function fmtDateTime(time, layout) {
    var ds = [...time.matchAll(/\d+/g)]
    if (ds.length < 3) return ""
    var p = function (idx) {
        var s = ds[idx][0]
        return (idx > 0) ? s.padStart(2, 0) : s
    }
    var s = layout.replace("Y", p(0)).replace("m", p(1)).replace("d", p(2))
    s = ds.length > 3 ? s.replace("H", p(3)) : s.replace("H", "00")
    s = ds.length > 4 ? s.replace("i", p(4)) : s.replace("i", "00")
    s = ds.length > 5 ? s.replace("s", p(5)) : s.replace("s", "00")
    return s
}