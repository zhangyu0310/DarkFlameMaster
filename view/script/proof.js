function doProof() {
    var proof = document.getElementById("proof").value;
    if (proof == "") {
        alert("输入有效交易单号！")
        return
    }
    document.getElementById("proofForm").submit();
}