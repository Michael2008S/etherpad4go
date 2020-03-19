var Pad = require("./Pad");
var revNum = 2;
console.log(Pad);
var apad  = new Pad.Pad("testpad");
console.log(apad);
var res = apad.getKeyRevisionNumber(revNum);
console.log(res,revNum);

console.log(Math.floor(51/100)*100);
