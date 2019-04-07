function toast(title, mesg) {
    $.toast({
        heading: '<b style="font-size:15px;font-family:courier">'+title+"</b>",
        text: '<span style="font-size:15px;font-family:courier">'+mesg+'</span>',
        position: 'mid-center',
        icon: 'warning',
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

$(document).ajaxError(function(_, xhr, settings) {
    //stat = xhr.status + " " + xhr.statusText
    //mesg = settings.type + ' ' + settings.url + '<br><p style="white-space:pre">' + xhr.responseText + '</p>'
    var stat = `操作失败`
    var mesg = `<p>${xhr.responseText}</p>`
    toast(stat, mesg)
});