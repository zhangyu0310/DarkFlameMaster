var canChooseNum = 0
var chooseNum = 0
var chooseSeat = new Map()
var chosenSeat = new Map()
var maxRow = 0
var maxCol = 0
var proof = ""
var order = ""
var beforeChild = null
var additional = ""

function doSubmit() {
    if (chooseNum == 0) {
        alert("请至少选择一个座位")
        return false
    }
    var pack = {"proof": proof, "seatInfo": [], "additional": additional}
    for (let [key, value] of chooseSeat) {
        pack["seatInfo"].push({
            "row": Number(value.getAttribute("row")),
            "col": Number(value.getAttribute("col"))
        })
    }
    document.getElementById("chooseData").setAttribute("value", JSON.stringify(pack))
    document.getElementById("chooseForm").submit();
    return true;
}

const btStatus = {
    btsHit: 0,
    btsNotHit: 1,
    btsSelected: 2,
    btsInvalid: 3,
    btsIndex: 4,
}

function makeSeatMsg() {
    var msg = new Array()
    for (let [key, value] of chosenSeat) {
        msg.push("第"+value.getAttribute("row")+"排，第"+value.getAttribute("col")+"座")
    }
    for (let [key, value] of chooseSeat) {
        msg.push("第"+value.getAttribute("row")+"排，第"+value.getAttribute("col")+"座")
    }
    return msg
}

function setSeatMsg() {
    var seatMsg = document.getElementById("seatMsg")
    while (seatMsg.hasChildNodes()) {
        seatMsg.removeChild(seatMsg.firstChild)
    }
    var msg = makeSeatMsg()
    for (let i = 0; i < msg.length; i++) {
        seatMsg.appendChild(document.createTextNode(msg[i]))
        seatMsg.appendChild(document.createElement("br"))
    }
}

function canChooseSeat() {
    return canChooseNum - chooseNum > 0
}

function clickSeat(id) {
    var bt = document.getElementById(id)
    const bts = bt.getAttribute("bts")
    if (bts == btStatus.btsNotHit) {
        if (!canChooseSeat()) {
            return
        }
        bt.style.backgroundColor = "#FDE6E0"
        bt.style.color = "#FDE6E0"
        bt.setAttribute("bts", btStatus.btsHit)
        chooseNum++
        chooseSeat.set(id, bt)
        setSeatMsg()
    } else if (bts == btStatus.btsHit) {
        bt.style.backgroundColor = "#C7EDCC"
        bt.style.color = "#C7EDCC"
        bt.setAttribute("bts", btStatus.btsNotHit)
        chooseNum--
        chooseSeat.delete(id)
        setSeatMsg()
    }
}

function hoverSeat(id) {
    var bt = document.getElementById(id)
    const bts = bt.getAttribute("bts")
    if (bts == btStatus.btsHit || bts == btStatus.btsNotHit) {
        if (!canChooseSeat() && bts == btStatus.btsNotHit) {
            return
        }
        bt.style.boxShadow = "0px 0px 10px #8A9585"
    }
}

function outSeat(id) {
    var bt = document.getElementById(id)
    const bts = bt.getAttribute("bts")
    if (bts == btStatus.btsHit || bts == btStatus.btsNotHit) {
        bt.style.boxShadow = ""
    }
}

function newRow() {
    return document.createElement('p')
}

function createRow(map) {
    var mapRow = newRow()
    map.appendChild(mapRow)
    mapRow.style.marginLeft = "0px"
    mapRow.style.marginRight = "0px"
    mapRow.style.marginTop = "6px"
    mapRow.style.marginBottom = "6px"
    return mapRow
}

function newSeat() {
    return document.createElement("button")
}

function createSeat(map, row, col, bts, isMine) {
    var bt = newSeat()
    if (order == "right") {
        map.insertBefore(bt, beforeChild)
        beforeChild = bt
    } else {
        map.appendChild(bt)
    }
    bt.style.width = "30px"
    bt.style.height = "30px"
    bt.style.marginLeft = "3px"
    bt.style.marginRight = "3px"
    bt.style.marginTop = "0px"
    bt.style.marginBottom = "0px"
    var id = row + "-" + col
    bt.setAttribute("id", id)
    bt.style.border = "none"
    bt.style.borderRadius = "30%"
    bt.style.transitionDelay = "0.1s"
    bt.style.transitionDuration = "0.2s"
    bt.setAttribute("bts", bts)
    bt.setAttribute("row", row)
    bt.setAttribute("col", col)
    bt.innerText = "·"
    switch (bts) {
        case btStatus.btsNotHit:
            bt.style.backgroundColor = "#C7EDCC";
            bt.style.color = "#C7EDCC"
            bt.style.cursor = "pointer"
            bt.setAttribute("onclick", "clickSeat(\"" + id + "\")")
            bt.setAttribute("onmouseover", "hoverSeat(\"" + id + "\")")
            bt.setAttribute("onmouseout", "outSeat(\"" + id + "\")")
            break;
        case btStatus.btsHit:
            bt.style.backgroundColor = "#FDE6E0";
            bt.style.color = "#FDE6E0"
            bt.style.cursor = "pointer"
            bt.setAttribute("onclick", "clickSeat(\"" + id + "\")")
            bt.setAttribute("onmouseover", "hoverSeat(\"" + id + "\")")
            bt.setAttribute("onmouseout", "outSeat(\"" + id + "\")")
            break;
        case btStatus.btsSelected:
            if (isMine) {
                bt.style.backgroundColor = "rgba(236,23,48,0.94)";
                bt.style.color = "rgba(236,23,48,0.94)"
                chosenSeat.set(id, bt)
            } else {
                bt.style.backgroundColor = "#8A9585";
                bt.style.color = "#8A9585"
            }
            bt.style.cursor = "not-allowed"
            break;
        case btStatus.btsInvalid:
            bt.style.backgroundColor = "transparent"
            bt.style.color = "transparent"
            bt.style.cursor = "default"
            break;
        case btStatus.btsIndex:
            bt.style.backgroundColor = "transparent"
            bt.style.cursor = "default"
            if (row == 0) {
                if (col == -1) {
                    bt.innerText = "·"
                    bt.style.color = "transparent"
                } else {
                    bt.innerText = (col+1).toString()
                }
            } else {
                bt.innerText = row.toString()
            }
            break;
    }
    bt.style.opacity = "0.85"
    return bt
}

function checkBrowser() {
    return function () {
        var ua = navigator.userAgent,
            isWindowsPhone = /(?:Windows Phone)/.test(ua),
            isSymbian = /(?:SymbianOS)/.test(ua) || isWindowsPhone,
            isAndroid = /(?:Android)/.test(ua),
            isFireFox = /(?:Firefox)/.test(ua),
            isChrome = /(?:Chrome|CriOS)/.test(ua),
            isTablet = /(?:iPad|PlayBook)/.test(ua) || (isAndroid && !/(?:Mobile)/.test(ua)) || (isFireFox && /(?:Tablet)/.test(ua)),
            isPhone = /(?:iPhone)/.test(ua) && !isTablet,
            isPc = !isPhone && !isAndroid && !isSymbian;
        return {
            isTablet: isTablet,
            isPhone: isPhone,
            isAndroid: isAndroid,
            isPc: isPc
        };
    }()
}

function printSi(si) {
    var msg = ""
    for (let i = 0; i < si.length; i++) {
        msg += si[i]["row"] + " " + si[i]["col"] + " " + si[i]["status"] + "\n"
    }
    return msg
}

function printBl(bl) {
    var msg = ""
    for (let i = 0; i < bl.length; i++) {
        msg += bl[i]["row"] + " " + bl[i]["col"] + " " + bl[i]["blockNum"] + "\n"
    }
    return msg
}

function seatSortFunc(a, b) {
    if (a["row"] > b["row"]) {
        return 1;
    } else if (a["row"] < b["row"]) {
        return -1;
    } else {
        if (a["col"] > b["col"]) {
            return 1;
        } else {
            return -1;
        }
    }
}

const blockDirection = {
    front: "front",
    back: "back",
}

function readData() {
    // 读取服务端传来的座位数据（JSON），并解析
    const data = document.getElementById("sd")
    // alert(data.value)
    var seatData = JSON.parse(data.value)
    if (seatData["msg"] != "") {
        alert(seatData["msg"])
    }
    canChooseNum = seatData["canChooseNum"]
    maxRow = seatData["maxRow"]
    maxCol = seatData["maxCol"]
    proof = seatData["proof"]
    order = seatData["order"]
    additional = seatData["additional"]
    // 读取排序座位信息
    let si = seatData["seatInfo"]
    si.sort(seatSortFunc)
    // alert(printSi(si))
    // 读取排序空格信息
    let bl = seatData["blockInfo"]
    bl.sort(seatSortFunc)
    // alert(printBl(bl))
    // 根据数据生成座位图
    var map = document.getElementById("seatMap")
    // 创建索引行 - 由于增加了空格，所以不需要创建索引行
    // var indexRow = createRow(map)
    // createSeat(indexRow, 0, -1, btStatus.btsIndex)
    // for (let i = 0; i < maxCol; i++) {
    //     createSeat(indexRow, 0, i, btStatus.btsIndex)
    // }
    // 创建座位行
    var blockIndex = 0
    var lastRow = 0;
    var mapRow;
    for (let i = 0; i < si.length; i++) {
        // 座位前的行号
        if (lastRow != si[i]["row"]) {
            mapRow = createRow(map)
            beforeChild = null
            createSeat(mapRow, si[i]["row"], 0, btStatus.btsIndex)
            lastRow = si[i]["row"]
        }
        // 匹配空格信息
        if (blockIndex < bl.length &&
            bl[blockIndex]["row"] == si[i]["row"] &&
            bl[blockIndex]["col"] == si[i]["col"] &&
            bl[blockIndex]["direction"] == blockDirection.front) {
            for (let j = 0; j < bl[blockIndex]["blockNum"]; j++) {
                createSeat(mapRow, -100-blockIndex, -j-1, btStatus.btsInvalid)
            }
            blockIndex++
        }
        createSeat(mapRow, si[i]["row"], si[i]["col"], si[i]["status"], si[i]["isMine"])
        // 匹配空格信息
        if (blockIndex < bl.length &&
            bl[blockIndex]["row"] == si[i]["row"] &&
            bl[blockIndex]["col"] == si[i]["col"] &&
            bl[blockIndex]["direction"] == blockDirection.back) {
            for (let j = 0; j < bl[blockIndex]["blockNum"]; j++) {
                createSeat(mapRow, -100-blockIndex, -j-1, btStatus.btsInvalid)
            }
            blockIndex++
        }
    }

    // 生成选座信息
    setSeatMsg()
}
