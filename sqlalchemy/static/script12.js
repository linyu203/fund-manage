//build-table.js
//
'use strict';

console.log("Build table init");

function generateTableHead(table,data) {
  let thead = table.createTHead();
  let row = thead.insertRow();
  console.log(data);
  for (let key of data) {
    let th = document.createElement("th");
    let text = document.createTextNode(key);
    th.appendChild(text);
    row.appendChild(th);
  }
}
function generateTable(table, data, th) {
  for (let element of data) {
    let row = table.insertRow();
    for (let idx in th) {
      let cell = row.insertCell();
      let text = document.createTextNode(element[th[idx]]);
      cell.appendChild(text);
    }
  }
}

function rowClicked(){
	console.log("Row was clicked!");
	//console.log(this);
	if(this.sectionRowIndex>0){
		var tds = $(this).find('td');
		let fundName = $(tds[0]).text();
		console.log(fundName);
	}

	
}

function InitRow(){
	//console.log("InitRow");
    $(document.body).on("click","tr", rowClicked);
	
}

const demo = [
	{Name: "demo", Count: 2, Description: "demo fund", Creation: "Jan 01 2020 10:30:20"}
]

function loadPage(funds){
    console.log(funds);
    let table = document.querySelector("table");
    let th = Object.keys(demo[0]);
    generateTableHead(table, th);
    if(funds !== undefined) {
	//let data = Object.keys(funds[0]);
	//generateTableHead(table, data);
	generateTable(table, funds, th);
    }else{
	//return 
        generateTable(table, demo, th);
    }
    InitRow();
}
function newfund(){
	var xhr = new XMLHttpRequest();
	xhr.open("CREATE","/",true);
	xhr.setRequestHeader("Content-Type","application/x-www-urlencoded");
	xhr.send("action=Create");
	console.log(xhr);
	
}
document.getElementById("RFR").addEventListener("click", function () {
	window.location.reload();
});
document.getElementById("CNF").addEventListener("click", function () {
	console.log("Create Fund clicked");
	//newfund();
	window.location.href = "fund";
});

