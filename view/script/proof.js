function doProof() {
    var proof = document.getElementById("proof").value;
    if (proof == "") {
        alert("输入为空！")
        return
    }
    document.getElementById("proofForm").submit();
}

function doCheck() {
    var proof = document.getElementById("proof").value;
    if (proof == "") {
        alert("输入为空！")
        return
    }
    document.getElementById("check").setAttribute("value", proof)
    document.getElementById("checkForm").submit();
}
