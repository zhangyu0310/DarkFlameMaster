var canChooseNum = 0
var chooseNum = 0
var chooseSeat = new Map()
var maxRow = 0
var maxCol = 0
var proof = ""

function doSubmit() {
    if (chooseNum == 0) {
        alert("请至少选择一个座位")
        return false
    }
    var pack = {"proof": proof, "seatInfo": []}
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

function createSeat(map, row, col, bts) {
    var bt = newSeat()
    map.appendChild(bt)
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
            bt.style.backgroundColor = "#8A9585";
            bt.style.color = "#8A9585"
            bt.style.cursor = "not-allowed"
            break;
        case btStatus.btsInvalid:
            bt.style.backgroundColor = "transparent"
            bt.style.color = "transparent"
            bt.style.cursor = "not-allowed"
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

function printSi(si) {
    var msg = ""
    for (let i = 0; i < si.length; i++) {
        msg += si[i]["row"] + " " + si[i]["col"] + " " + si[i]["status"] + "\n"
    }
    return msg
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
    // 根据数据生成座位图
    let si = seatData["seatInfo"]
    si.sort(function (a, b) {
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
    })
    // alert(printSi(si))
    var map = document.getElementById("seatMap")
    // 创建索引行
    var indexRow = createRow(map)
    createSeat(indexRow, 0, -1, btStatus.btsIndex)
    for (let i = 0; i < maxCol; i++) {
        createSeat(indexRow, 0, i, btStatus.btsIndex)
    }
    // 创建座位行
    var mapRow;
    for (let i = 0; i < si.length; i++) {
        // 座位前的行号
        if (i % maxCol == 0) {
            mapRow = createRow(map)
            createSeat(mapRow, si[i]["row"], 0, btStatus.btsIndex)
        }
        createSeat(mapRow, si[i]["row"], si[i]["col"], si[i]["status"])
    }

    // 生成选座信息
    setSeatMsg()
}
