var proofName = ""
var additionalName = ""

function doProof() {
    var proof = document.getElementById("proof").value;
    if (proof == "") {
        alert("["+proofName+"] 输入为空！请确认全部填写。")
        return
    }
    if (additionalName != "") {
        var additional = document.getElementById("additional").value;
        if (additional == "") {
            alert("["+additionalName+"] 输入为空！请确认全部填写。")
            return
        }
    }
    document.getElementById("proofForm").submit();
}

function doCheck() {
    var proof = document.getElementById("proof").value;
    if (proof == "") {
        alert("校验信息输入为空！")
        return
    }
    document.getElementById("check").setAttribute("value", proof)
    document.getElementById("checkForm").submit();
}

function loadLabel() {
    proofName = document.getElementById("proofInput").value
    additionalName = document.getElementById("additionalInput").value
    if (additionalName == "") {
        document.getElementById("proofForm").removeChild(document.getElementById("additionalPart"))
    }
    var proofPart = document.getElementById("proofPart")
    var pf = proofPart.firstChild
    proofPart.insertBefore(document.createTextNode(proofName), pf)
    proofPart.insertBefore(document.createElement("br"), pf)
    if (additionalName != "") {
        var additionalPart = document.getElementById("additionalPart")
        var af = additionalPart.firstChild
        additionalPart.insertBefore(document.createTextNode(additionalName), af)
        additionalPart.insertBefore(document.createElement("br"), af)
    }
}
