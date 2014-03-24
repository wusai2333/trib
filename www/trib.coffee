me = ""
showing = ""

listTribs = (data) ->
    ret = JSON.parse(data)
    if ret.Err != ""
        appendError(ret.Err)
        return
 
    tribs = $("div#tribs")
    tribs.empty()

    if ret.Tribs.length == 0
        tribs.append("No Tribble.")
        return

    ul = $("<ul/>")
    ret.Tribs.reverse()
    for trib in ret.Tribs
        li = $("<li/>")
        li.append('<span class="author">@' + trib.User + '</span> ')
        li.append('<span class="time">' + trib.Time + '</span> ')
        li.append($('<span class="trib" />').text(trib.Message))
        ul.append(li)
    tribs.append(ul)

    return

showHome = (ev) ->
    ev.preventDefault()
    _showHome()
    return

_showHome = ->
    console.log("show home: " + me)
    $.ajax({
        url: "api/list-home"
        type: "POST"
        data: me
        success: listTribs
        cache: false
    })

    showing = "!home"
    
    $("div#timeline").show()
    $("div#whom").hide()
    $("a#follow").hide()
    $("div#tribs").empty()
    $("h2#title").html("Home of " + me)

    return

showUser = (ev) ->
    ev.preventDefault()
    name = $(this).text()
    console.log("show user: " + name)
    $.ajax({
        url: "api/list-tribs"
        type: "POST"
        data: name
        success: listTribs
        cache: false
    })

    showing = name
    $("h2#title").html(name)

    $("div#tribs").empty()
    $("div#timeline").show()
    $("div#whom").show()
    $("a#follow").show()
    updateFollow()

    return

updateUsers = (data) ->
    ret = JSON.parse(data)
    if ret.Err != ""
        appendError(ret.Err)
        return

    users = $("#users")
    users.empty()
    if ret.Users.length == 0
        users.append("No user.")
        return

    ul = $("<ul/>")
    for name in ret.Users
        ul.append('<li><a href="#">' + 
            name + '</a></li>')
    users.append(ul)
    $("#users li").click(showUser)

    return
    
addUser = ->
    name = $("form#adduser input#username").val()
    if name == ""
        return false

    $("form#adduser input#username").val("")

    console.log("add user", name)
    $.ajax({
        url: "api/add-user"
        type: "POST"
        data: name
        success: updateUsers
        cache: false
    })
    
    return false

listUsers = ->
    $.ajax({
        url: "api/list-users"
        success: updateUsers
        cache: false
    })
    return

appendError = (e) ->
    $("div#errors").show()
    $("div#errors").append('<div class="error">Error: ' +
        e + '</div>')

postTrib = ->
    #TODO
    return false

signIn = (ev) ->
    ev.preventDefault()
    if showing == "" || showing == "!home"
        return

    console.log("sigin in as: " + showing)

    me = showing
    $("div#who").show()
    $("div#who h3").html("Signed in as " + me)
    $("div#compose").show()

    _showHome()
    updateFollow()

    return

signOut = (ev) ->
    console.log("sign out")

    ev.preventDefault()
    me = ""
    $("div#who").hide()
    $("div#compose").hide()
    $("a#follow").hide()

    if showing == "!home"
        $("div#timeline").hide()

    return

hoveringFollow = false

_updateFollow = (data) ->
    but = $("a#follow")
    ret = JSON.parse(data)
    if ret.Err != ""
        appendError(ret.Err)
        return

    but.unbind("mouseenter mouseleave")
    but.unbind("click")
    if ret.V
        if hoveringFollow
            but.html("Unfollow")
        else
            but.html("Following")
        but.hover(((ev) ->
            $(this).html("Unfollow")
            hoveringFollow = true
            return
        ), ((ev) ->
            $(this).html("Following")
            hoveringFollow = false
            return
        ))
        but.click(unfollow)
    else
        but.html("Follow")
        but.hover(((ev) ->
            hoveringFollow = true
            return
        ), ((ev) ->
            hoveringFollow = false
            return
        ))
        but.click(follow)

    return

follow = (ev) ->
    ev.preventDefault()
    $.ajax({
        url: "api/follow"
        type: "POST"
        data: JSON.stringify({
            Who: me
            Whom: showing
        })
        success: _updateFollow
        cache: false
    })
    return

unfollow = (ev) ->
    ev.preventDefault()
    $.ajax({
        url: "api/unfollow"
        type: "POST"
        data: JSON.stringify({
            Who: me
            Whom: showing
        })
        success: _updateFollow
        cache: false
    })
    return

updateFollow = ->
    if me == "" || showing == "!home"
        $("a#follow").hide()
        return

    $("a#follow").html("&nbsp;")
    $.ajax({
        url: "api/is-following"
        type: "POST"
        data: JSON.stringify({
            Who: me
            Whom: showing
        })
        success: _updateFollow
        cache: false
    })
    return

main = ->
    $("form#adduser").submit(addUser)
    $("form#post").submit(postTrib)

    $("div#errors").hide()
    $("div#timeline").hide()

    $("a#signin").click(signIn)
    $("a#home").click(showHome)
    $("a#signout").click(signOut)

    listUsers()
    return

$(document).ready(main)

