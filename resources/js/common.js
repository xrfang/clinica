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