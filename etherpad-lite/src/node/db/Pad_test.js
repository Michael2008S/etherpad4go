var Pad = require("./Pad");
var revNum = 2;
console.log(Pad);
var apad = new Pad.Pad("testpad");
console.log(apad);
var res = apad.getKeyRevisionNumber(revNum);
console.log(res, revNum);

console.log(Math.floor(51 / 100) * 100);



var Changeset = require("ep_etherpad-lite/static/js/Changeset");

text = "Welcome to Etherpad!\\n\\nThis pad text is synchronized~ https:\\/\\/github.com\\/ether\\/etherpad-lite\\n"
let cleanTxt = Pad.cleanText(text);
let firstChangeset = Changeset.makeSplice("\n", 0, 0, cleanTxt);
console.log(firstChangeset);
