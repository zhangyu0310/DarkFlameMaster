var order = "right"
var beforeChild = null

const btStatus = {
    btsHit: 0,
    btsNotHit: 1,
    btsIndex: 2,
}

function makeSeat() {
    const maxRow = document.getElementById("maxRow").value
    const maxCol = document.getElementById("maxCol").value
    order = document.getElementById("order").value

    var map = document.getElementById("seatMap")
    var mapRow;
    for (let i = 0; i < maxRow; i++) {
        beforeChild = null
        mapRow = createRow(map)
        createSeat(mapRow, i+1, 0, btStatus.btsIndex)
        for (let j = 0; j < maxCol; j++) {
            createSeat(mapRow, i, j, btStatus.btsNotHit)
        }
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
